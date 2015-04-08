package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/keybase/client/go/engine"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

// CmdDeviceAdd is the 'device add' command.  It is used for
// device provisioning to enter a secret phrase on an existing
// device.
type CmdDeviceAdd struct {
	phrase string
}

// NewCmdDeviceAdd creates a new cli.Command.
func NewCmdDeviceAdd(cl *libcmdline.CommandLine) cli.Command {
	return cli.Command{
		Name:        "add",
		Usage:       "keybase device add \"secret phrase\"",
		Description: "Authorize a new device",
		Action: func(c *cli.Context) {
			cl.ChooseCommand(&CmdDeviceAdd{}, "add", c)
		},
	}
}

// RunClient runs the command in client/server mode.
func (c *CmdDeviceAdd) RunClient() error {
	cli, err := GetDeviceClient()
	if err != nil {
		return err
	}
	protocols := []rpc2.Protocol{
		NewSecretUIProtocol(),
		NewDoctorUIProtocol(),
	}
	if err := RegisterProtocols(protocols); err != nil {
		return err
	}

	return cli.DeviceAdd(c.phrase)
}

// Run runs the command in standalone mode.
func (c *CmdDeviceAdd) Run() error {
	ctx := &engine.Context{SecretUI: G_UI.GetSecretUI(), DoctorUI: G_UI.GetDoctorUI()}
	eng := engine.NewKexSib(G, c.phrase)
	return engine.RunEngine(eng, ctx)
}

// ParseArgv gets the secret phrase from the command args.
func (c *CmdDeviceAdd) ParseArgv(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return fmt.Errorf("device add takes one arg: the secret phrase")
	}
	c.phrase = ctx.Args()[0]
	return nil
}

// GetUsage says what this command needs to operate.
func (c *CmdDeviceAdd) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:    true,
		KbKeyring: true,
		API:       true,
	}
}
