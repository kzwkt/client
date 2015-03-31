//
//  KBSecretWordsView.m
//  Keybase
//
//  Created by Gabriel on 3/17/15.
//  Copyright (c) 2015 Gabriel Handford. All rights reserved.
//

#import "KBSecretWordsView.h"

@interface KBSecretWordsView ()
@property KBLabel *label;
@property KBLabel *secretWordsLabel;
@end

@implementation KBSecretWordsView

- (void)viewInit {
  [super viewInit];
  self.backgroundColor = KBAppearance.currentAppearance.backgroundColor;

  YOView *contentView = [[YOView alloc] init];
  [self addSubview:contentView];

  KBLabel *header = [[KBLabel alloc] init];
  [header setText:@"Register Device" style:KBTextStyleHeaderLarge alignment:NSCenterTextAlignment lineBreakMode:NSLineBreakByTruncatingTail];
  [contentView addSubview:header];

  _label = [[KBLabel alloc] init];
  [contentView addSubview:_label];

  _secretWordsLabel = [[KBLabel alloc] init];
  _secretWordsLabel.selectable = YES;
  [_secretWordsLabel setBackgroundColor:KBAppearance.currentAppearance.secondaryBackgroundColor];
  [_secretWordsLabel setBorderEnabled:YES];
  _secretWordsLabel.insets = UIEdgeInsetsMake(10, 20, 10, 20);
  [contentView addSubview:_secretWordsLabel];

  YOView *footerView = [YOView view];
  _button = [KBButton buttonWithText:@"OK" style:KBButtonStylePrimary];
  [footerView addSubview:_button];
  _cancelButton = [KBButton buttonWithText:@"Cancel" style:KBButtonStyleDefault];
  [footerView addSubview:_cancelButton];
  footerView.viewLayout = [YOLayout layoutWithLayoutBlock:[KBLayouts layoutForButton:_button cancelButton:_cancelButton horizontalAlignment:KBHorizontalAlignmentCenter]];
  [contentView addSubview:footerView];

  YOSelf yself = self;
  contentView.viewLayout = [YOLayout layoutWithLayoutBlock:^(id<YOLayout> layout, CGSize size) {
    CGFloat y = 0;

    y += [layout centerWithSize:CGSizeMake(400, 0) frame:CGRectMake(40, y, size.width - 80, 0) view:header].size.height + 20;

    y += [layout centerWithSize:CGSizeMake(400, 0) frame:CGRectMake(40, y, size.width - 80, 0) view:yself.label].size.height + 30;

    y += [layout centerWithSize:CGSizeMake(500, 0) frame:CGRectMake(40, y, size.width - 80, 0) view:yself.secretWordsLabel].size.height + 40;

    y += [layout sizeToFitVerticalInFrame:CGRectMake(40, y, size.width - 80, 0) view:footerView].size.height;

    return CGSizeMake(MIN(580, size.width), y);
  }];

  self.viewLayout = [YOLayout layoutWithLayoutBlock:[KBLayouts center:contentView]];
}

- (void)viewDidAppear:(BOOL)animated {
  [self.window recalculateKeyViewLoop];
  [self.window makeFirstResponder:_secretWordsLabel.textView];
}

- (void)setSecretWords:(NSString *)secretWords deviceNameToRegister:(NSString *)deviceNameToRegister {
  /*
   On your "CLI" computer, a window should have appeared. Type this in it:

       term unfold yellow reflect spare now

   Alternatively, if you're using the terminal at "CLI", type this:

       keybase sibkey add "term unfold yellow reflect spare now"
   */

  // NSStringWithFormat(@"In order to register this device you need to enter in these secret words on the device named: <strong>%@</strong>.", deviceName)
  [_label setMarkup:@"In order to register this device you need to enter in the secret phrase generated on an existing device." style:KBTextStyleDefault alignment:NSLeftTextAlignment lineBreakMode:NSLineBreakByWordWrapping];
  [_secretWordsLabel setText:secretWords font:[NSFont fontWithName:@"Monaco" size:20] color:KBAppearance.currentAppearance.textColor alignment:NSCenterTextAlignment lineBreakMode:NSLineBreakByWordWrapping];
  [self setNeedsLayout];
}

@end
