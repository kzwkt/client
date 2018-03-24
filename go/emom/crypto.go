package emom

import (
	binary "encoding/binary"
	emom1 "github.com/keybase/client/go/protocol/emom1"
	clockwork "github.com/keybase/clockwork"
	saltpack "github.com/keybase/saltpack"
	context "golang.org/x/net/context"
	sync "sync"
)

func makeNonce(msgType emom1.MsgType, n emom1.Seqno) saltpack.Nonce {
	var out saltpack.Nonce
	copy(out[0:16], "encrypted_fmprpc")
	binary.BigEndian.PutUint32(out[12:16], uint32(msgType))
	binary.BigEndian.PutUint64(out[16:], uint64(n))
	return out
}

func encrypt(ctx context.Context, msg []byte, msgType emom1.MsgType, n emom1.Seqno, key saltpack.BoxPrecomputedSharedKey) (emom1.AuthEnc, error) {
	if key == nil {
		return emom1.AuthEnc{}, NoSessionKeyError
	}
	return emom1.AuthEnc{
		N: n,
		E: key.Box(makeNonce(msgType, n), msg),
	}, nil
}

func decrypt(ctx context.Context, msgType emom1.MsgType, ae emom1.AuthEnc, key saltpack.BoxPrecomputedSharedKey) ([]byte, error) {
	return key.Unbox(makeNonce(msgType, ae.N), ae.E)
}

type ServerPublicKey struct {
	gen emom1.KeyGen
	key saltpack.BoxPublicKey
}

type User struct {
	uid            emom1.UID
	userSigningKey saltpack.SigningSecretKey
}

type UsersCryptoPackage struct {
	user            User
	serverPublicKey ServerPublicKey
	sessionKey      saltpack.BoxPrecomputedSharedKey
	clock           clockwork.Clock
}

func (u *UsersCryptoPackage) InitClient(ctx context.Context, arg *emom1.Arg, rp *emom1.RequestPlaintext) error {

	if u.sessionKey == nil {
		return nil
	}

	ephemeralKey, err := u.serverPublicKey.key.CreateEphemeralKey()
	if err != nil {
		return err
	}
	ephemeralKID := emom1.KID(ephemeralKey.GetPublicKey().ToKID())

	u.sessionKey = ephemeralKey.Precompute(u.serverPublicKey.key)

	handshake := emom1.Handshake{
		V: 1,
		S: u.serverPublicKey.gen,
		K: ephemeralKID,
	}

	authToken := emom1.AuthToken{
		C: emom1.ToTime(u.clock.Now()),
		D: emom1.KID(u.user.userSigningKey.GetPublicKey().ToKID()),
		K: ephemeralKID,
		U: u.user.uid,
	}

	msg, err := encodeToBytes(authToken)
	if err != nil {
		return err
	}
	sig, err := u.user.userSigningKey.Sign(msg)
	if err != nil {
		return err
	}
	signedAuthToken := emom1.SignedAuthToken{
		T: authToken.Export(),
		S: sig,
	}

	arg.H = &handshake
	rp.F = &signedAuthToken

	return nil
}

func (u *UsersCryptoPackage) SessionKey() saltpack.BoxPrecomputedSharedKey {
	return u.sessionKey
}

func (u *UsersCryptoPackage) InitServerHandshake(_ context.Context, _ emom1.Arg) error {
	return nil
}

func (u *UsersCryptoPackage) InitUserAuth(_ context.Context, _ emom1.Arg, _ emom1.RequestPlaintext) error {
	return nil
}

func (c *UsersCryptoPackage) ServerRatchet(ctx context.Context, res *emom1.Res) error {
	return nil
}

var _ Cryptoer = (*UsersCryptoPackage)(nil)

type CloudCryptoPackage struct {
	sync.Mutex
	serverKeys           map[emom1.KeyGen]saltpack.BoxSecretKey
	userAuth             func(context.Context, emom1.UID, emom1.KID) error
	importSigningKey     func(context.Context, emom1.KID) (saltpack.SigningPublicKey, error)
	checkReplayAndImport func(context.Context, emom1.KID) (saltpack.BoxPublicKey, error)
	user                 emom1.UID
	clock                clockwork.Clock
	ratchetSeqno         emom1.Seqno
	serverKeyInUse       saltpack.BoxSecretKey
	userKeyInUse         saltpack.BoxPublicKey

	// SessionKeys. Seqno=0 is with long-lived server public key. Seqno>0 are with
	// ratcheted server keys, which can later be thrown away.
	sessionKeys map[emom1.Seqno]saltpack.BoxPrecomputedSharedKey
}

func (c *CloudCryptoPackage) SessionKey() saltpack.BoxPrecomputedSharedKey {
	c.Lock()
	defer c.Unlock()
	return c.sessionKey()
}

func (c *CloudCryptoPackage) sessionKey() saltpack.BoxPrecomputedSharedKey {
	key, _ := c.sessionKeys[c.ratchetSeqno]
	return key
}

func (c *CloudCryptoPackage) setMasterSessionKey(k saltpack.BoxPrecomputedSharedKey) {
	c.sessionKeys[emom1.Seqno(0)] = k
}

func (c *CloudCryptoPackage) InitServerHandshake(ctx context.Context, arg emom1.Arg) error {
	c.Lock()
	defer c.Unlock()
	if c.sessionKey() == nil {
		return nil
	}
	if arg.H == nil {
		return NewHandshakeError("expected a handshake, but none given")
	}
	if arg.H.V != 1 {
		return NewHandshakeError("Can only support V1, got %d", arg.H.V)
	}
	userEphemeralKey, err := c.checkReplayAndImport(ctx, arg.H.K)
	if err != nil {
		return err
	}

	key, found := c.serverKeys[arg.H.S]
	if !found {
		return NewHandshakeError("key generation %d not found", arg.H.S)
	}

	c.serverKeyInUse = key
	c.userKeyInUse = userEphemeralKey
	c.setMasterSessionKey(key.Precompute(userEphemeralKey))

	return nil
}

func (c *CloudCryptoPackage) ServerRatchet(ctx context.Context, res *emom1.Res) error {
	c.Lock()
	defer c.Unlock()
	if c.ratchetSeqno > emom1.Seqno(0) {
		return nil
	}

	nextEphemeralKey, err := c.serverKeyInUse.GetPublicKey().CreateEphemeralKey()
	if err != nil {
		return err
	}
	nextSessionKey := nextEphemeralKey.Precompute(c.userKeyInUse)
	nextSeqno := c.ratchetSeqno + 1

	ratchet := emom1.ServerRatchet{
		I: nextSeqno,
		K: emom1.KID(nextEphemeralKey.GetPublicKey().ToKID()),
	}
	var encodedRatchet []byte
	var encryptedRatchet emom1.AuthEnc

	encodedRatchet, err = encodeToBytes(ratchet)
	if err != nil {
		return err
	}
	encryptedRatchet, err = encrypt(ctx, encodedRatchet, emom1.MsgType_RATCHET, nextSeqno, c.sessionKey())
	res.R = &encryptedRatchet
	c.ratchetSeqno = nextSeqno
	c.sessionKeys[nextSeqno] = nextSessionKey

	return nil
}

func (c *CloudCryptoPackage) InitClient(ctx context.Context, arg *emom1.Arg, rp *emom1.RequestPlaintext) error {
	return nil
}

func (c *CloudCryptoPackage) InitUserAuth(ctx context.Context, arg emom1.Arg, rp emom1.RequestPlaintext) error {

	// The user is authed and there's no more auth information coming down. Perfact!
	if c.user != nil && rp.F == nil {
		return nil
	}

	if c.user != nil && rp.F != nil {
		return newUserAuthError("attempt to reauth an already-authed session")
	}

	if arg.H == nil {
		return newUserAuthError("User auth must happen along with the handshake")
	}

	at := emom1.AuthToken{
		C: rp.F.T.C,
		K: arg.H.K,
		U: rp.F.T.U,
	}

	encodedAuthToken, err := encodeToBytes(at)
	if err != nil {
		return err
	}

	userKey, err := c.importSigningKey(ctx, rp.F.T.D)
	if err != nil {
		return err
	}

	err = userKey.Verify(encodedAuthToken, rp.F.S)
	if err != nil {
		return err
	}

	err = c.userAuth(ctx, rp.F.T.U, rp.F.T.D)
	if err != nil {
		return err
	}

	return nil
}

var _ Cryptoer = (*CloudCryptoPackage)(nil)
