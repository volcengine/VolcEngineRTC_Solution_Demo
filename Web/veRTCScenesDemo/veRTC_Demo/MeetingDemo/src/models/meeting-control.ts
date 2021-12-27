/**
 * @author fuyuhao
 */

import {actionCreatorFactory} from 'dva-model-creator';
import { AppModel } from '@/app-interfaces';
import {setFieldReducer} from '@/utils/redux-utils';
import MeetingControlSDK from '@/lib/MeetingController';

export interface MeetingControlModelState {
  sdk: Nullable<MeetingControlSDK>,
}

export const meetingControlInitialState: MeetingControlModelState = {
  sdk: null,
};

const factory = actionCreatorFactory('meetingControl');

export const meetingControlActions = {
  initSDK: factory<MeetingControlSDK>('initSDK'),
};

const LoginModel: AppModel<MeetingControlModelState> = {
  namespace: 'meetingControl',
  state: meetingControlInitialState,
  subscriptions: {
    setup({dispatch}) {
      dispatch(meetingControlActions.initSDK(new MeetingControlSDK()));
    },
  },
  reducers: {
    initSDK: setFieldReducer(meetingControlInitialState, 'sdk'),
  }
};


export default LoginModel;
