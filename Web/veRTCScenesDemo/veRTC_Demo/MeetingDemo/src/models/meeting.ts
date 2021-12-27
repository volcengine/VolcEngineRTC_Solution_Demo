import {actionCreatorFactory} from 'dva-model-creator';
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
  meetingInfo: MeetingInfo;
  meetingUsers: MeetingUser[];
  orderMeetingUsers: MeetingUser[];

  status: 'end' | 'start' | 'init' | 'closeTips' | 'hidden' | 'lockTrackEnded';
  speakCollapse: boolean;
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

  speakCollapse: false,

  meetingInfo: {
    created_at: Date.now(),
    host_id: '',
    now: Date.now(),
    record: false,
    room_id: '',
    screen_shared_uid: '',
  },
  meetingUsers: [],
  orderMeetingUsers: [],
  status: 'init'
};

const factory = actionCreatorFactory('meeting');

export const meetingActions = {
  setMeetingInfo: factory<MeetingModelState['meetingInfo']>('setMeetingInfo'),
  setMeetingUsers:
    factory<MeetingModelState['meetingUsers']>('setMeetingUsers'),
  setMeetingOrderUsers: factory<MeetingModelState['meetingUsers']>(
    'setMeetingOrderUsers'
  ),
  setViewMode: factory<MeetingModelState['viewMode']>('setViewMode'),
  setMeetingStatus: factory<MeetingModelState['status']>('setMeetingStatus'),
  setSpeakCollapse:
    factory<MeetingModelState['speakCollapse']>('setSpeakCollapse'),
};

const MeetingModel: AppModel<MeetingModelState> = {
  namespace: 'meeting',
  state: meetingInitialState,
  reducers: {
    setMeetingInfo: setFieldReducer(meetingInitialState, 'meetingInfo'),
    setMeetingUsers: setFieldReducer(meetingInitialState, 'meetingUsers'),
    setMeetingOrderUsers: setFieldReducer(
      meetingInitialState,
      'orderMeetingUsers'
    ),
    setViewMode: setFieldReducer(meetingInitialState, 'viewMode'),
    setMeetingStatus: setFieldReducer(meetingInitialState, 'status'),
    setSpeakCollapse: setFieldReducer(meetingInitialState, 'speakCollapse'),
  },
};
export default MeetingModel;
