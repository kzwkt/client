//
//  KBLaunchCtl.h
//  Keybase
//
//  Created by Gabriel on 3/12/15.
//  Copyright (c) 2015 Gabriel Handford. All rights reserved.
//

#import <Foundation/Foundation.h>

#import <GHKit/GHKit.h>
#import "KBDefines.h"

typedef void (^KBLaunchExecution)(NSError *error, NSString *output);

@interface KBLaunchCtl : NSObject

@property BOOL releaseOnly;

- (void)load:(KBLaunchExecution)completion;
- (void)unload:(KBLaunchExecution)completion;

- (void)reload:(KBLaunchExecution)completion;

- (void)status:(KBLaunchExecution)completion;

- (void)installLaunchAgent:(void (^)(NSError *error))completion;

@end
