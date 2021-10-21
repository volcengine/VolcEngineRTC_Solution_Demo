#import "MeetingRTCManager.h"
#import "RoomVideoSession.h"
#import "RoomViewController.h"
#import <Foundation/Foundation.h>
#import <ByteRtcEngineKit/ByteRtcEngineKit.h>
#import "RoomVideoSession.h"
#import "RoomParamInfoModel.h"

NS_ASSUME_NONNULL_BEGIN
@class MeetingRTCManager;
@protocol MeetingRTCManagerDelegate <NSObject>

- (void)meetingRTCManager:(MeetingRTCManager *)meetingRTCManager changeParamInfo:(RoomParamInfoModel *)model;

- (void)rtcManager:(MeetingRTCManager * _Nonnull)rtcManager didStreamAdded:(NSString *_Nullable)streamsUid;

- (void)rtcManager:(MeetingRTCManager * _Nonnull)rtcManager didScreenStreamAdded:(NSString *_Nullable)screenStreamsUid;

- (void)rtcManager:(MeetingRTCManager *_Nonnull)rtcManager didStreamRemoved:(NSString *_Nullable)streamsUid;

- (void)rtcManager:(MeetingRTCManager *_Nonnull)rtcManager didScreenStreamRemoved:(NSString *)screenStreamsUid;

- (void)rtcManager:(MeetingRTCManager *_Nonnull)rtcManager reportAllAudioVolume:(NSDictionary<NSString *, NSNumber *> *_Nonnull)volumeInfo;

@end

@interface MeetingRTCManager : NSObject

@property (nonatomic, weak) id<MeetingRTCManagerDelegate> delegate;

/*
 * RTC Manager Singletons
 */
+ (MeetingRTCManager *_Nullable)shareRtc;

#pragma mark - Base Method

/**
 * Create RTCEngine instance
 * @param appID The unique identifier of each application is randomly generated by the VRTC console. Instances generated by different AppIds are completely independent for audio and video calls in VRTC and cannot communicate with each other.
 */
- (void)createEngine:(NSString *)appID;

/**
 * Join room
 * @param videoSession User Model
 */
- (void)joinChannelWithRoomVideoSession:(RoomVideoSession *)videoSession;

/*
 * Real-time update of video parameters
 */
- (void)updateRtcVideoParams;

/*
 * Switch camera
 */
- (void)switchCamera;

/*
 * Switch audio routing (handset/speaker)
 * @param enableSpeaker ture:Use speakers false：Use the handset
 */
- (int)setEnableSpeakerphone:(BOOL)enableSpeaker;

/*
 * Switch local audio capture
 * @param mute ture:Turn on audio capture false：Turn off audio capture
 */
- (void)enableLocalAudio:(BOOL)enable;

/*
 * Switch local video capture
 * @param mute ture:Open video capture false：Turn off video capture
 */
- (void)enableLocalVideo:(BOOL)enable;

/*
 * Leave the room
 */
- (void)leaveChannel;

/*
 * destroy
 */
- (void)destroy;

/*
 * Open preview
 @param view View
 */
- (void)startPreview:(UIView *)view;

/*
 * gGet Sdk Version
 */
- (NSString *_Nullable)getSdkVersion;

#pragma mark - Subscribe Stream

/*
 * Subscribe to the video stream
 @param uid User ID
 */
- (void)subscribeStream:(NSString *_Nullable)uid;

/*
 * Subscribe to screen stream
 @param uid User ID
 */
- (void)subscribeScreenStream:(NSString *)uid;

/*
 * Unsubscribe video stream
 @param uid User ID
 */
- (void)unsubscribe:(NSString *_Nullable)uid;

/*
 * Unsubscribe from screen stream
 @param uid User ID
 */
- (void)unsubscribeScreen:(NSString *)uid;

/*
 * Remote render view and uid binding
 @param canvas Canvas Model
 */
- (void)setupRemoteVideo:(ByteRtcVideoCanvas *)canvas;

/*
 * Local render view and uid binding
 @param canvas Canvas Model
 */
- (void)setupLocalVideo:(ByteRtcVideoCanvas *)canvas;

#pragma mark - Screen

/*
 * Screen render view and uid binding
 * @param view View
 * @param uid User ID
 */
- (void)setupRemoteScreen:(UIView *)view uid:(NSString *)uid;


/*
 * Turn on screen sharing
 * @param param Param data
 * @param extension Extension Bunle ID
 * @param groupId Group ID
 */
- (void)startScreenSharingWithParam:(ScreenCaptureParam *_Nonnull)param preferredExtension:(NSString *_Nullable)extension groupId:(NSString *_Nonnull)groupId;

/*
 * Turn off screen sharing
 */
- (void)stopScreenSharing;


/*
 * Get current subscription Uid
 * @return subscription Uid
 */
- (NSDictionary *_Nullable)getSubscribeUidDic;

@end

NS_ASSUME_NONNULL_END