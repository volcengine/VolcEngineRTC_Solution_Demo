import { UserModelState } from './models/user';
import { MeetingModelState } from './models/meeting';
import { MeetingSettingsState } from './models/meeting-settings';
import { MeetingControlModelState } from '@@/plugin-dva/connect';
import {ImmerReducer} from '@@/plugin-dva/connect';
import {Action} from 'dva-model-creator';
import {EffectsMapObject, SubscriptionsMapObject} from 'dva';
import { AudioStats } from '@/lib/socket-interfaces';
export interface AppState {
  user: UserModelState,
  meeting: MeetingModelState,
  meetingControl: MeetingControlModelState,
  meetingSettings: MeetingSettingsState,
}

export interface AppModel<S> {
  namespace: string;
  state: S;
  reducers: {
    [K in string]: ImmerReducer<S, Action<any>>
  },
  effects?: EffectsMapObject,
  subscriptions?: SubscriptionsMapObject,
}

export enum FitType {
  cover, contain
}

type EncoderConfiguration = {
  resolution?: {
    width: number;
    height: number;
  };
  frameRate?: {
    min: number;
    max: number;
  };
  bitrate?: {
    min: number;
    max: number;
  };
};

export type IStreamStats = {
  accessDelay: number;
  audioSentBytes: number;
  audioReceivedBytes: number;
  videoSentResolutionWidth: number;
  videoSentResolutionHeight: number;
  videoSentFrameRate: number;
  videoSentBitrate: number;
  audioSentBitrate: number;
  videoReceivedResolutionWidth: number;
  videoReceivedResolutionHeight: number;
  videoReceivedPacketsLost: number;
  audioReceivedPacketsLost: number;
};

export interface Stream {
  uid: string;
  stream: {
    screen: boolean;
  };
  getId: () => string;
  enableAudio: () => void;
  disableAudio: () => void;
  enableVideo: () => void;
  disableVideo: () => void;
  close: () => void;
  init: (onSuccess: () => void, onFailure: (err: Error) => void) => void;
  play: (id: string, options?: { fit?: FitType; muted?: boolean }) => void;
  setVideoEncoderConfiguration: (options: EncoderConfiguration) => void;
  getStats(): IStreamStats;
  getAudioLevel(): number;
  playerComp: JSX.Element;
}

type SubscribeOption = {
  video?: boolean;
  audio?: boolean;
  bigStream?: boolean;
  forceSubscribe?: boolean;
};

export interface RTCClint {
  init: (
    appId: string,
    onSuccess?: () => void,
    onFailure?: (err: Error) => void
  ) => void;
  join: (
    token: string,
    roomId: string,
    userId: string,
    onSuccess?: (uid: string) => void,
    onFailure?: (err: Error) => void
  ) => void;
  getConnectionState: () => string;
  publish: (stream: Stream) => Promise<void>;
  unpublish: (
    stream: Stream,
    onFailure?: (err: Error) => void
  ) => Promise<void>;
  subscribe: (
    stream: Stream,
    options: SubscribeOption,
    onFailure?: (err: Error) => void
  ) => void;
  leave: (onSuccess?: () => void, onFailure?: (err: Error) => void) => void;
  on: (event: string, callback: (param: any) => void) => void;
  getRemoteAudioStats(callback: (param: AudioStats) => void): void;
  getLocalAudioStats(callback: (param: AudioStats) => void): void;
}

export interface ConnectStatus {
  connected: boolean;
  disconnect: boolean;
}

export interface DeviceInstance {
  deviceId: string;
  groupId: string;
  kind: 'audioinput' | 'audiooutput' | 'videoinput';
  label: string;
}

export type DeviceItems = {
  audioinput : DeviceInstance[];
  videoinput : DeviceInstance[];
  audiooutput : DeviceInstance[];
}

export type IRemoteAudioLevel = { userId: string | null; RecvLevel: number };
