/**
 * @author yuanyuan
 */

import { AppModel } from '@/app-interfaces';
import RtcClient from '@/rtcApi/rtc-client';

export interface RTCClientControlModelState {
  rtc: RtcClient;
}

export const rtcClientInitialState: RTCClientControlModelState = {
  rtc: new RtcClient(),
};

const LoginModel: AppModel<RTCClientControlModelState> = {
  namespace: 'rtcClientControl',
  state: rtcClientInitialState,
  subscriptions: {},
  reducers: {},
};

export default LoginModel;
