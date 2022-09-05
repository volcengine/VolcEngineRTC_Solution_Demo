import { IVolume } from './../app-interfaces';
import { actionCreatorFactory } from 'dva-model-creator';
import { AppModel } from '@/app-interfaces';
import { setFieldReducer } from '@/utils/redux-utils';

/**
 * @author fuyuhao
 */

export enum ViewMode {
  /**
   * default mode
   */
  GalleryView,
  /**
   * Switch to this mode when someone shares the screen
   * 当有人分享屏幕时会切换到此模式
   */
  SpeakerView,
}
export interface MeetingInfo {
  created_at: number;
  host_id: string;
  now: number;
  /**
   * Recording
   * 是否正在录制
   */
  record: boolean;
  room_id: string;
  /**
   * The ID of the user currently sharing the screen
   * 当前正在屏幕共享的用户ID
   */
  screen_shared_uid: string;
  // hasTrackError: { kind: string; isScreen: boolean } | null;
  localSpeaker: IVolume;
  volumeSortList: IVolume[];
}

export interface MeetingUser {
  created_at: number;
  is_camera_on: boolean;
  is_host: boolean;
  is_mic_on: boolean;
  is_sharing: boolean;
  room_id: string;
  user_id: string;
  user_name: string;
  user_uniform_id: string;
}

export interface MeetingModelState {
  appId: string | null;
  token: string | null;
  roomId: string | null;
  userId: string | null;

  resolution: string;
  frameRate: number;
  bitrate: number;

  audioInputs: MediaDeviceInfo[];
  videoInputs: MediaDeviceInfo[];
  curMic: MediaDeviceInfo | null;
  curCam: MediaDeviceInfo | null;

  screenSharing: boolean;
  screenResolution: string;
  screenBitrate: number;
  screenFramerate: number;
  showStatus: boolean;

  viewMode: ViewMode;
  localCaptureSuccess: boolean;
  localAudioVideoCaptureSuccess: boolean;
  meetingInfo: MeetingInfo;
  meetingUsers: MeetingUser[];

  status: 'end' | 'start' | 'init' | 'closeTips' | 'hidden' | 'lockTrackEnded';
  speakCollapse: boolean;
  isSharing: boolean;
}

const meetingInitialState: MeetingModelState = {
  appId: null,
  token: null,
  roomId: null,
  userId: null,

  resolution: '1920*1080',
  frameRate: 15,
  bitrate: 1000,

  audioInputs: [],
  videoInputs: [],
  curMic: null,
  curCam: null,

  screenSharing: false,
  screenResolution: '1920*1080',
  screenBitrate: 1000,
  screenFramerate: 15,
  showStatus: false,

  viewMode: ViewMode.GalleryView,

  localCaptureSuccess: false,
  localAudioVideoCaptureSuccess: false,
  speakCollapse: false,

  meetingInfo: {
    created_at: Date.now(),
    host_id: '',
    now: Date.now(),
    record: false,
    room_id: '',
    screen_shared_uid: '',
    localSpeaker: {
      volume: 0,
      userId: '',
    },
    volumeSortList: [],
    // hasTrackError: null,
  },
  meetingUsers: [],
  status: 'init',
  isSharing: false,
};

const factory = actionCreatorFactory('meeting');

export const meetingActions = {
  setMeetingInfo: factory<MeetingModelState['meetingInfo']>('setMeetingInfo'),
  setMeetingUsers:
    factory<MeetingModelState['meetingUsers']>('setMeetingUsers'),
  setViewMode: factory<MeetingModelState['viewMode']>('setViewMode'),
  setMeetingStatus: factory<MeetingModelState['status']>('setMeetingStatus'),
  setSpeakCollapse:
    factory<MeetingModelState['speakCollapse']>('setSpeakCollapse'),
  setLocalCaptureSuccess: factory<MeetingModelState['localCaptureSuccess']>(
    'setLocalCaptureSuccess'
  ),
  setLocalAudioVideoCaptureSuccess: factory<
    MeetingModelState['localAudioVideoCaptureSuccess']
  >('setLocalAudioVideoCaptureSuccess'),
  setMeetingIsSharing: factory<MeetingModelState['isSharing']>(
    'setMeetingIsSharing'
  ),
};

const MeetingModel: AppModel<MeetingModelState> = {
  namespace: 'meeting',
  state: meetingInitialState,
  reducers: {
    setMeetingInfo: setFieldReducer(meetingInitialState, 'meetingInfo'),
    setMeetingUsers: setFieldReducer(meetingInitialState, 'meetingUsers'),
    setViewMode: setFieldReducer(meetingInitialState, 'viewMode'),
    setMeetingStatus: setFieldReducer(meetingInitialState, 'status'),
    setSpeakCollapse: setFieldReducer(meetingInitialState, 'speakCollapse'),
    setLocalCaptureSuccess: setFieldReducer(
      meetingInitialState,
      'localCaptureSuccess'
    ),
    setLocalAudioVideoCaptureSuccess: setFieldReducer(
      meetingInitialState,
      'localAudioVideoCaptureSuccess'
    ),
    setMeetingIsSharing: setFieldReducer(meetingInitialState, 'isSharing'),
  },
};
export default MeetingModel;
