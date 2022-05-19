import { actionCreatorFactory } from 'dva-model-creator';
import { AppModel } from '@/app-interfaces';
import { setFieldReducer, setFieldsReducer } from '@/utils/redux-utils';
import { TOASTS } from '@/config';
/**
 * @author fuyuhao
 */

export interface UserModelState {
  appId: string | null;
  name: string | null;
  roomId: string | null;
  isHost: boolean;
  isMicOn: boolean;
  isCameraOn: boolean;
  isSharing: boolean;
  createdAt: number;
  logged: boolean;
  userId: string | null;
  deviceAccess: {
    audio: boolean;
    video: boolean;
    audioMessage?: keyof typeof TOASTS;
    videoMessage?: keyof typeof TOASTS;
  };
  network: boolean;
}

export const userInitialState: UserModelState = {
  appId: null,
  name: null,
  roomId: null,
  isHost: false,
  isMicOn: true,
  isCameraOn: true,
  isSharing: false,
  createdAt: -1,
  logged: false,
  userId: null,
  deviceAccess: {
    audio: true,
    video: true,
  },
  network: true,
};

const factory = actionCreatorFactory('user');

export const userActions = {
  setAppId: factory<UserModelState['appId']>('setAppId'),
  setUserName: factory<UserModelState['name']>('setUserName'),
  setUserId: factory<UserModelState['name']>('setUserId'),
  setRoomId: factory<UserModelState['roomId']>('setRoomId'),
  setIsMicOn: factory<UserModelState['isMicOn']>('setIsMicOn'),
  setIsCameraOn: factory<UserModelState['isCameraOn']>('setIsCameraOn'),
  setIsSharing: factory<UserModelState['isSharing']>('setIsSharing'),
  setCreateAt: factory<UserModelState['createdAt']>('setCreateAt'),
  setIsHost: factory<UserModelState['isHost']>('setIsHost'),
  setUserFields: factory<Partial<UserModelState>>('setUserFields'),
  setLogged: factory<UserModelState['logged']>('setLogged'),
  setDeviceAccess: factory<UserModelState['deviceAccess']>('setDeviceAccess'),
  setNetWork: factory<UserModelState['network']>('setNetWork'),
};

const UserModel: AppModel<UserModelState> = {
  namespace: 'user',
  state: userInitialState,
  subscriptions: {},
  reducers: {
    setAppId: setFieldReducer(userInitialState, 'appId'),
    setUserName: setFieldReducer(userInitialState, 'name'),
    setUserId: setFieldReducer(userInitialState, 'userId'),
    setRoomId: setFieldReducer(userInitialState, 'roomId'),
    setIsMicOn: setFieldReducer(userInitialState, 'isMicOn'),
    setIsCameraOn: setFieldReducer(userInitialState, 'isCameraOn'),
    setIsSharing: setFieldReducer(userInitialState, 'isSharing'),
    setCreateAt: setFieldReducer(userInitialState, 'createdAt'),
    setIsHost: setFieldReducer(userInitialState, 'isHost'),
    setUserFields: setFieldsReducer<UserModelState>(),
    setLogged: setFieldReducer(userInitialState, 'logged'),
    setDeviceAccess: setFieldReducer(userInitialState, 'deviceAccess'),
    setNetWork: setFieldReducer(userInitialState, 'network'),
  },
};

export default UserModel;
