import io from 'socket.io-client';
import { EventEmitter } from 'eventemitter3';
import { v4 as uuid } from 'uuid';
import { SOCKETURL, SOCKETPATH, TOASTS } from '@/config';
import Utils from '@/utils/utils';
import Logger from '@/utils/Logger';
import type { MeetingInfo, MeetingUser } from '@/models/meeting';
import { Stream, RTCClint, ConnectStatus } from '@/app-interfaces';
import VRTC from '@volcengine/rtc';

import type {
  SocketResponse,
  GetAppIDResponse,
  JoinMeetingPayload,
  JoinMeetingResponse,
  UserPayload,
  GetAppIDPayload,
  RecordMeetingPayload,
  GetVerifyCodePayload,
  VerifyLoginSms,
  VerifyLoginRes,
  VerifyLoginToken,
  EndMeetingPayload,
  HistoryVideoRecord,
  AudioStats,
} from './socket-interfaces';

const logger = new Logger('MettingController');
type Socket = typeof io.Socket;

type IStreams = {
  [key: string]: Stream;
};

type SubscribeOption = {
  video?: boolean;
  audio?: boolean;
};

type SubscribePatch = {
  subscribed: string[];
  unsubscribed: string[];
};
class MettingController extends EventEmitter {
  private _isHost: boolean;
  private socket!: Socket;
  private client!: RTCClint;
  private reconnect!: boolean;
  private eventListenNames: string[] = [
    'onUserMicStatusChange',
    'onUserCameraStatusChange',
    'onHostChange',
    'onUserJoinMeeting',
    'onUserLeaveMeeting',
    'onShareScreenStatusChanged',
    'onRecord',
    'onMeetingEnd',
    'onMuteAll',
    'onMuteUser',
    'onAskingMicOn',
    'onUserKickedOff',
  ];

  constructor() {
    super();
    logger.debug('MettingController constructor()');

    this._isHost = false;
    this.reconnect = false;
    window.__mc = this;
  }

  public connect(): Promise<string> {
    return new Promise((resolve, reject) => {
      logger.debug('connect()');

      const options = {
        secure: true,
        transports: ['websocket'],
        query: {
          wsid: uuid(),
          appid: 'veRTCDemo',
          ua: `web-${VRTC.getSdkVersion()}`,
          did: Utils.getDeviceId(),
        },
        path: SOCKETPATH,
      };

      if (this.socket && this.socket.disconnected) {
        this.socket.connect();
      } else {
        this.socket = io.connect(SOCKETURL, options);
      }

      this.socket.on('connect', () => {
        this._handleSocket();

        if (this.reconnect) {
          this.userReconnect().then((res) => {
            this.emit('onUserReconnect', res);
          });
          this.reconnect = false;
        }

        // RTC.Logger.setLogLevel(RTC.Logger.DEBUG);
        // logger.debug('client: %o', this.client);

        resolve(this.socket.id);
      });

      this.socket.on('connect_error', reject);

      this.socket.on('reconnect', () => {
        this.reconnect = true;
      });

      this.socket.on('reconnecting', () => {
        //alert('reconnecting');
      });

      this.socket.on('reconnect_failed', () => {
        // alert('reconnect_failed')
      });
      this.socket.on('disconnect', () => {
        // alert('disconnect');
      });
    });
  }

  public disconnect(): void {
    if (this.socket?.connected) {
      //this.socket.disconnect();
    }
  }

  public getPhoneVerifyCode(payload: GetVerifyCodePayload): Promise<string> {
    return this.sendSignaling('sendLoginSms', payload);
  }

  public verifyLoginSms(payload: VerifyLoginSms): Promise<VerifyLoginRes> {
    return this.sendSignaling('verifyLoginSms', payload);
  }

  public verifyLoginToken(payload: VerifyLoginToken): Promise<string | null> {
    return this.sendSignaling('verifyLoginToken', payload);
  }

  public getAppID(payload: GetAppIDPayload): Promise<GetAppIDResponse> {
    return this.sendSignaling<GetAppIDPayload, GetAppIDResponse>(
      'getAppID',
      payload
    );
  }

  public getHistoryVideoRecord(): Promise<HistoryVideoRecord[]> {
    return this.sendSignaling('getHistoryVideoRecord');
  }

  public deleteVideoRecord(payload: { vid: string }): Promise<null> {
    return this.sendSignaling('deleteVideoRecord', payload);
  }

  public joinMeeting(p: JoinMeetingPayload): Promise<JoinMeetingResponse> {
    return new Promise((resolve, reject) => {
      this.sendSignaling<JoinMeetingPayload, JoinMeetingResponse>(
        'joinMeeting',
        p
      )
        .then((response) => {
          resolve(response);
        })
        .catch(reject);
    });
  }

  public userReconnect(): Promise<number> {
    return new Promise((resolve, reject) => {
      if (!this.socket.connected) {
        reject('websocket disconnected');
        return;
      }
      const login_token = Utils.getLoginToken();
      this.socket.emit(
        'userReconnect',
        {
          login_token,
        },
        (soscketRet: SocketResponse<{ code: number }>) => {
          resolve(soscketRet.code);
        }
      );
    });
  }

  public leaveMeeting(payload: VerifyLoginToken): Promise<null> {
    return new Promise((resolve, reject) => {
      // this.clientLeave(
      // () => {
      this.sendSignaling('leaveMeeting', payload).finally(() => {
        this.disconnect();
        resolve(null);
      });
      // },
      // () => {
      //   reject('Leave RTC Room failed');
      // }
      // );
    });
  }

  public turnOnMic(payload: VerifyLoginToken): Promise<null> {
    return this.sendSignaling('turnOnMic', payload);
  }

  public turnOffMic(payload: VerifyLoginToken): Promise<null> {
    return this.sendSignaling('turnOffMic', payload);
  }

  public turnOnCamera(payload: VerifyLoginToken): Promise<null> {
    return this.sendSignaling('turnOnCamera', payload);
  }

  public turnOffCamera(payload: VerifyLoginToken): Promise<null> {
    return this.sendSignaling('turnOffCamera', payload);
  }

  public getMeetingUserInfo(payload: UserPayload): Promise<MeetingUser[]> {
    return this.sendSignaling('getMeetingUserInfo', payload);
  }

  public getMeetingInfo(): Promise<MeetingInfo> {
    return this.sendSignaling<null, MeetingInfo>('getMeetingInfo');
  }

  public startShareScreen(payload: VerifyLoginToken): Promise<null> {
    return this.sendSignaling('startShareScreen', payload);
  }

  public endShareScreen(): Promise<null> {
    return this.sendSignaling('endShareScreen');
  }

  public changeHost(payload: UserPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('changeHost', payload);
  }

  public muteUser(payload: UserPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('muteUser', payload);
  }

  public askMicOn(payload: UserPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('askMicOn', payload);
  }

  public askCameraOn(payload: UserPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('askCameraOn', payload);
  }

  public endMeeting(payload: EndMeetingPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('endMeeting', payload);
  }

  public recordMeeting(payload: RecordMeetingPayload): Promise<null> {
    try {
      this._assertNotHost();
    } catch (error) {
      throw TOASTS['record'];
    }
    return this.sendSignaling('recordMeeting', payload);
  }

  public updateRecordLayout(payload: RecordMeetingPayload): Promise<null> {
    this._assertNotHost();
    return this.sendSignaling('updateRecordLayout', payload);
  }

  public sendSignaling<P, T>(type: string, payload?: P): Promise<T> {
    return new Promise((resolve, reject) => {
      if (!this.socket.connected) {
        reject('websocket disconnected');
        return;
      }

      const login_token = Utils.getLoginToken();
      const body = {
        ...payload,
        login_token,
      };

      this.socket.emit(type, body, (soscketRet: SocketResponse<T>) => {
        if (soscketRet.code === 200) {
          resolve(soscketRet.response);
        } else {
          reject(soscketRet.message);
        }
      });
    });
  }

  // public publish(stream: Stream): void {
  //   this._assertNotInRoom();
  //   this.client.publish(stream);
  // }

  public unpublish(stream: Stream): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this._assertNotInRoom();
        this.client.unpublish(stream).then(() => resolve());
      } catch (error) {
        reject();
      }
    });
  }

  public getConnectStatus(): ConnectStatus {
    if (this.socket) {
      return {
        connected: this.socket.connected,
        disconnect: this.socket.disconnected,
      };
    }
    return {
      connected: false,
      disconnect: false,
    };
  }

  public checkSocket(): Promise<string | undefined> {
    const socketStatus = this.getConnectStatus();
    if (socketStatus.connected) {
      logger.debug('已连接');
      return Promise.resolve();
    } else {
      logger.debug('未连接');
      return new Promise((resolve, reject) => {
        this?.connect()
          .then((socketId) => {
            logger.debug('mc connect successfully', socketId);
            resolve(socketId);
          })
          .catch((err) => reject(err));
      });
    }
  }

  private _handleSocket() {
    logger.debug('_handleSocket()');

    this.eventListenNames.forEach((type) => {
      this.socket.on(type, (payload: any) => {
        this.emit(type, payload.data);
      });
    });
  }

  public removeEvent(): void {
    [...this.eventListenNames, 'onUserReconnect'].forEach((type) => {
      this.removeListener(type);
    });
  }

  private _assertNotHost() {
    if (!this._isHost) {
      throw new Error('Permission Denied');
    }
  }

  private _assertNotInRoom() {
    if (!this.client || this.client.getConnectionState() !== 'CONNECTED') {
      throw new Error('Not in RTC Room');
    }
  }

  set isHost(v: boolean) {
    this._isHost = v;
  }
}

export default MettingController;
