//
//  MeetingControlTool.m
//  SceneRTCDemo
//
//  Created by  on 2021/3/17.
//

#import "MeetingControlTool.h"

@implementation MeetingControlTool

+ (NSString *)convertToJsonData:(NSDictionary *)dict {
    if (!dict) {
        return @"";
    }
    NSError *error;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:dict options:NSJSONWritingPrettyPrinted error:&error];
    NSString *jsonString;
    if (!jsonData) {
        NSLog(@"%@", error);
    } else {
        jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
    }
    NSMutableString *mutStr = [NSMutableString stringWithString:jsonString];
    NSRange range = {0, jsonString.length};
    [mutStr replaceOccurrencesOfString:@" " withString:@"" options:NSLiteralSearch range:range];
    NSRange range2 = {0, mutStr.length};
    [mutStr replaceOccurrencesOfString:@"\n" withString:@"" options:NSLiteralSearch range:range2];
    return mutStr;
}

+ (MeetingControlAckModel *)dataToAckModel:(NSArray *)dataList {
    MeetingControlAckModel *model = nil;
    if (dataList && [dataList isKindOfClass:[NSArray class]] && dataList.count > 0) {
        NSDictionary *dic = dataList.firstObject;
        if (dic && [dic isKindOfClass:[NSDictionary class]] && dic.count > 0) {
            model = [[MeetingControlAckModel alloc] init];
            model.code = [dic[@"code"] integerValue];
            model.response = dic[@"response"];
            NSString *message = [MeetingControlTool messageToRequestCode:model.code];
            if (!message || ![message isKindOfClass:[NSString class]] || message.length <= 0) {
                message = dic[@"message"];
            }
            model.message = message;
            if (model.code == 450) {
                [[NSNotificationCenter defaultCenter] postNotificationName:NotificationLoginExpired object:nil];
            }
        }
    }
    return model;
}

+ (NSString *)messageToRequestCode:(NSInteger)code {
    NSString *message = @"";
    switch (code) {
        case 406:
            message = @"?????????????????????";
            break;
        case 412:
            message = @"??????????????????????????????????????????????????????????????????";
            break;
        case 422:
            message = @"??????????????????";
            break;
        case 430:
            message = @"?????????????????????????????????????????????";
            break;
        case 440:
            message = @"??????????????????????????????????????????";
            break;
        case 441:
            message = @"?????????????????????????????????????????????";
            break;
        case 450:
            message = @"????????????????????????????????????";
            break;
        case 503:
            message = @"??????????????????????????????";
            break;
        case 504:
            message = @"?????????????????????????????????";
            break;
        case 702:
            message = @"?????????Token????????????????????????";
            break;
        default:
            break;
    }
    return message;
}

+ (MeetingControlNoticeModel *)dataToNoticeModel:(NSArray *)dataList {
    MeetingControlNoticeModel *model = nil;
    if (dataList && [dataList isKindOfClass:[NSArray class]] && dataList.count > 0) {
        NSDictionary *dic = dataList.firstObject;
        if (dic && [dic isKindOfClass:[NSDictionary class]] && dic.count > 0) {
            model = [[MeetingControlNoticeModel alloc] init];
            model.data = dic[@"data"];
        }
    }
    return model;
}

+ (NSDictionary *)addToken:(NSDictionary *)dic {
    NSMutableDictionary *tokenDic = nil;
    if (dic && [dic isKindOfClass:[NSDictionary class]] && dic.count > 0) {
        tokenDic = [dic mutableCopy];
        [tokenDic setValue:[MenuTokenCompoments token] forKey:@"login_token"];
    } else {
        tokenDic = [[NSMutableDictionary alloc] init];
        [tokenDic setValue:[MenuTokenCompoments token] forKey:@"login_token"];
    }
    return [tokenDic copy];
}

@end
