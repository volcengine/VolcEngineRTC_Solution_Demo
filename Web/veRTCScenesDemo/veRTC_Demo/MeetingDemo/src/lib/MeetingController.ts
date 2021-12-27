import RTC from '@/sdk/VRTC.esm.min.js';
import io from 'socket.io-client';
import { EventEmitter } from 'eventemitter3';
import { v4 as uuid } from 'uuid';
import { SOCKETURL, SOCKETPATH } from '@/config';
import Utils from '@/utils/utils';
import Logger from '@/utils/Logger';
import type { MeetingInfo, MeetingUser } from '@/models/meeting';
import { Stream, RTCClint, ConnectStatus } from '@/app-interfaces';
import { TOASTS } from '@/config';

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
  private count = 8; //订阅远端流数量 8 + 1（本地流） = 9
  private _streams: IStreams = {};
  private streamRecord: SubscribePatch = { subscribed: [], unsubscribed: [] };
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
          ua: `web-${RTC.version}`,
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
        if (!this.client) {
          this.client = RTC.createClient({
            iceUrl: process.env.ICEURL,
          });
          this._handleSocket();
          this._handleRTCEvents();
        }

        if (this.reconnect) {
          this.userReconnect().then((res) => {
            this.emit('onUserReconnect', res);
          });
          this.reconnect = false;
        }

        RTC.Logger.setLogLevel(RTC.Logger.DEBUG);
        logger.debug('client: %o', this.client);

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
          this.client.init(p.app_id, () => {
            this.client.join(
              response.token,
              p.room_id,
              p.user_id,
              (uid: string) => {
                logger.debug('RTC Room uid: ' + uid);
                resolve(response);
              },
              (err: unknown) => {
                logger.error(err);
                reject(err);
              }
            );
          });
          return response;
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
      this.clientLeave(
        () => {
          this.sendSignaling('leaveMeeting', payload).finally(() => {
            this.disconnect();
            resolve(null);
          });
        },
        () => {
          reject('Leave RTC Room failed');
        }
      );
    });
  }

  public clientLeave(success: () => void, fail?: () => void): void {
    this.client.leave(success, fail);
  }

  public getRemoteAudioStats(): Promise<AudioStats> {
    return new Promise((resolve) => {
      this.client.getRemoteAudioStats((param) => {
        resolve(param);
      });
    });
  }

  public getLocalAudioStats(): Promise<AudioStats> {
    return new Promise((resolve) => {
      this.client.getLocalAudioStats((param) => {
        resolve(param);
      });
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

  public publish(stream: Stream): void {
    this._assertNotInRoom();
    this.client.publish(stream);
  }

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

  public subscribe(stream: Stream, option?: SubscribeOption): void {
    this._assertNotInRoom();
    this.client.subscribe(
      stream,
      Object.assign(option || {}, { forceSubscribe: true })
    );
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

  public checkSocket(): Promise<string | void> {
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

  public changSubscribe(newUserSort: MeetingUser[]): void {
    newUserSort.forEach((i, index) => {
      if (!this._streams[i.user_id] || !i.user_id) {
        return;
      }
      const subscribe_index = this.streamRecord?.subscribed?.findIndex(
        (item) => item === i.user_id
      );
      const unsubscribe_index = this.streamRecord?.unsubscribed?.findIndex(
        (item) => item === i.user_id
      );

      if (index <= this.count) {
        if (subscribe_index === -1) {
          this.streamRecord?.subscribed?.push(i.user_id);
          this.subscribe(this._streams[i.user_id], {
            audio: true,
            video: true,
          });
        }
        if (unsubscribe_index !== -1) {
          this.streamRecord?.unsubscribed?.splice(unsubscribe_index, 1);
        }
      }
      if (index > this.count) {
        if (unsubscribe_index === -1) {
          this.streamRecord?.unsubscribed?.push(i.user_id);
          this.subscribe(this._streams[i.user_id], {
            audio: true,
            video: false,
          });
        }
        if (subscribe_index !== -1) {
          this.streamRecord?.subscribed?.splice(subscribe_index, 1);
        }
      }
    });
  }

  public cleanStreamRecord(): void {
    this.streamRecord = { subscribed: [], unsubscribed: [] };
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

  private _handleRTCEvents() {
    logger.debug('_handleRTCEvents()');

    this.client.on('stream-added', (payload: { stream: Stream }) => {
      logger.debug('stream-added: %o', payload);
      this.subscribe(payload.stream);
    });

    this.client.on('stream-subscribed', (payload: { stream: Stream }) => {
      logger.debug('stream-subscribed: %o', payload);
      this.emit('OnReceivedStream', payload.stream);
    });

    this.client.on('stream-removed', (payload: { stream: Stream }) => {
      logger.debug('stream-removed: %o', payload);
      this.emit('onRemoveStream', {
        uid: payload.stream.getId(),
        screen: payload.stream.stream.screen,
      });
    });

    // this.client.on('peer-leave', (payload: { uid: string }) => {
    //   logger.debug('stream-removed: %o', payload);
    //   this.emit('onUserLeaveMeetingClient', payload.uid);
    // });
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

  set streams(v: IStreams) {
    this._streams = v;
  }

  set isHost(v: boolean) {
    this._isHost = v;
  }
}

export default MettingController;
