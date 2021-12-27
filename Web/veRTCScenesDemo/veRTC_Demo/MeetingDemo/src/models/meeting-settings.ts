import { actionCreatorFactory } from 'dva-model-creator';
import { setFieldReducer } from '@/utils/redux-utils';
import { AppModel } from '@/app-interfaces';

interface Resolution {
  width: number;
  height: number;
}

interface Range {
  min: number;
  max: number;
}

interface StreamSettings {
  resolution: Resolution;
  frameRate: Range;
  bitrate: Range;
}

export interface MeetingSettingsState {
  streamSettings: StreamSettings;
  screenStreamSettings: StreamSettings;
  mic: string;
  camera: string;
  realtimeParam: boolean;
}

const streamInitialState = {
  resolution: {
    width: 640,
    height: 480,
  },
  frameRate: {
    min: 10,
    max: 15,
  },
  bitrate: {
    min: 250,
    max: 600,
  },
};

const screenStreamInitialState = {
  resolution: {
    width: 1920,
    height: 1080,
  },
  frameRate: {
    min: 10,
    max: 15,
  },
  bitrate: {
    min: 800,
    max: 2000,
  },
};

export const meetingSettingsInitialState: MeetingSettingsState = {
  streamSettings: streamInitialState,
  screenStreamSettings: screenStreamInitialState,
  mic: '',
  camera: '',
  realtimeParam: false,
};

const factory = actionCreatorFactory('meetingSettings');

export const meetingSettingsActions = {
  setStreamSettings: factory<MeetingSettingsState['streamSettings']>('setStreamSettings'),
  setScreenStreamSettings: factory<MeetingSettingsState['screenStreamSettings']>('setScreenStreamSettings'),
  setMic: factory<MeetingSettingsState['mic']>('setMic'),
  setCamera: factory<MeetingSettingsState['camera']>('setCamera'),
  setRealtimeParam: factory<MeetingSettingsState['realtimeParam']>('setRealtimeParam'),
};


const MeetingSettingsModel: AppModel<MeetingSettingsState> = {
  namespace: 'meetingSettings',
  state: meetingSettingsInitialState,
  subscriptions: {},
  reducers: {
    setStreamSettings: setFieldReducer(meetingSettingsInitialState, 'streamSettings'),
    setScreenStreamSettings: setFieldReducer(meetingSettingsInitialState, 'screenStreamSettings'),
    setMic: setFieldReducer(meetingSettingsInitialState, 'mic'),
    setCamera: setFieldReducer(meetingSettingsInitialState, 'camera'),
    setRealtimeParam: setFieldReducer(meetingSettingsInitialState, 'realtimeParam'),
  }
};

export default MeetingSettingsModel;
