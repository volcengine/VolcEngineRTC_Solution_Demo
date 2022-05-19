import VERTC, {
  SubscribeMediaType,
  MuteState,
  VideoRenderMode,
  IRTCEngine,
  RTCStream,
  LocalStreamStats,
  RemoteStreamStats,
  StreamIndex,
  RTCDevice
} from '@volcengine/rtc';

import { IVolume } from '@/app-interfaces';

export interface IBindEvent {
  handleStreamAdd: (params: { stream: RTCStream }) => void;
  handleStreamRemove: (params: { stream: RTCStream }) => void;
  handleEventError: (e: any, VERTC: any) => void;
  handleAudioVolumeIndication: (params: { speakers: IVolume[] }) => void;
  handleLocalStreamState: (stats: LocalStreamStats) => void;
  handleRemoteStreamState: (stats: RemoteStreamStats) => void;
  handleTrackEnd: (params: { kind: string; isScreen: boolean }) => void;
}

export default class RtcClient {
  config!: { appId: string; uid: string };
  engine!: IRTCEngine;
  handleStreamAdd!: (params: { stream: RTCStream }) => void;
  handleStreamRemove!: (params: { stream: RTCStream }) => void;

  constructor() {
    this.setRemoteVideoPlayer = this.setRemoteVideoPlayer.bind(this);
    this.setLocalVideoPlayer = this.setLocalVideoPlayer.bind(this);
    this.createScreenStream = this.createScreenStream.bind(this);
  }
  SDKVERSION = VERTC.getSdkVersion();
  init(props: { config: { appId: string; uid: string } }): void {
    this.config = props.config;
  }
  bindEngineEvents({
    handleStreamAdd,
    handleStreamRemove,
    handleEventError,
    handleAudioVolumeIndication,
    handleLocalStreamState,
    handleRemoteStreamState,
    handleTrackEnd,
  }: IBindEvent): void {
    this.handleStreamAdd = handleStreamAdd;
    this.engine.on(VERTC.events.onStreamAdd, handleStreamAdd);
    this.engine.on(VERTC.events.onStreamRemove, handleStreamRemove);
    this.engine.on(VERTC.events.onError, (e) => handleEventError(e, VERTC));
    this.engine.on(
      VERTC.events.onAudioVolumeIndication,
      handleAudioVolumeIndication
    );
    this.engine.on(VERTC.events.onLocalStreamStats, handleLocalStreamState);
    this.engine.on(VERTC.events.onRemoteStreamStats, handleRemoteStreamState);
    this.engine.on(VERTC.events.onTrackEnded, handleTrackEnd);
  }
  /**
   * remove the listeners when `createengine`
   */
  removeEventListener(): void {
    this.engine.off(VERTC.events.onStreamAdd, this.handleStreamAdd);
    this.engine.off(VERTC.events.onStreamRemove, this.handleStreamRemove);
  }
  join(token: string, roomId: string, uid: string): Promise<void> {
    return this.engine.joinRoom(
      token,
      roomId,
      {
        userId: uid,
      },
      {
        // 默认值全为false
        isAutoPublish: false,
        isAutoSubscribeAudio: false,
        isAutoSubscribeVideo: false,
      }
    );
  }
  // check permission of browser
  checkPermission(): Promise<{
    video: boolean;
    audio: boolean;
  }> {
    return VERTC.enableDevices();
  }
  /**
   * get the devices
   * @returns
   */
  async getDevices(): Promise<{
    audioInputs: RTCDevice[];
    videoInputs: RTCDevice[];
    audioPlaybackList: RTCDevice[];
  }> {
    return {
      audioInputs: await VERTC.getMicrophones(),
      videoInputs: await VERTC.getCameras(),
      audioPlaybackList: await VERTC.getAudioPlayback(),
    };
  }

  /**
   * @brief 取消共享屏幕音视频流
   * @function destoryScreenStream
   */
  async destoryScreenStream(
    success: () => void,
    fail: (err: Error) => void
  ): Promise<void> {
    //  停止捕获屏幕流
    this.engine
      .stopScreenCapture()
      .then(() => {
        this.engine.unpublishScreen();
        success && success();
      })
      .catch((err: Error) => fail && fail(err));
  }

  /**
   * @brief 创建本地音视频流
   * @function createLocalStream
   * @param streamOptions 流参数
   * @param isPublish 是否发布 预览则不发布
   */
  async createLocalStream(
    streamOptions: {
      mic: string;
      camera: string;
      audio: boolean;
      video: boolean;
    },
    isPublish: boolean,
    callback: (param: any) => void
  ): Promise<void> {
    const { mic, camera, audio, video } = streamOptions;
    const devicesStatus = {
      video: 1,
      audio: 1,
    };
    const permissions = await this.checkPermission();
    const devices = await this.getDevices();
    if (audio && permissions.audio) {
      await this.engine.startAudioCapture(
        mic ? mic : devices.audioInputs[0].deviceId
      );
    } else {
      if (!permissions.audio) devicesStatus['audio'] = 0;
      this.engine.muteLocalAudio(MuteState.MUTE_STATE_OFF);
    }
    if (video && permissions.video) {
      await this.engine.startVideoCapture(
        camera ? camera : devices.videoInputs[0].deviceId
      );
    } else {
      if (!permissions.video) devicesStatus['video'] = 0;
      this.engine.muteLocalVideo(MuteState.MUTE_STATE_OFF);
    }

    // 如果joinRoom的config设置了自动发布，这里就不需要发布了
    isPublish && (await this.engine.publish());

    callback &&
      callback({
        code: 0,
        msg: '设备获取成功',
        devicesStatus,
      });
  }

  /**
   * @brief 挂载流播放的容器
   * @param type 流类型 0/1
   */
  setLocalVideoPlayer(
    type: StreamIndex,
    renderDom: string | HTMLElement
  ): void {
    this.engine.setLocalVideoPlayer(type, {
      renderDom: renderDom,
      userId: this.config.uid,
      renderMode: VideoRenderMode.RENDER_MODE_HIDDEN,
    });
  }
  /**
   * @brief 挂载流播放的容器
   * @param type 流类型 0/1
   */
  async setRemoteVideoPlayer(
    type: StreamIndex,
    remoteUserId: string,
    domId: string | HTMLElement
  ): Promise<void> {
    this.engine
      .subscribeUserStream(
        remoteUserId,
        SubscribeMediaType.AUDIO_AND_VIDEO,
        type
      )
      .then(() => {
        try {
          this.engine.setRemoteVideoPlayer(type, {
            userId: remoteUserId,
            renderDom: domId,
            renderMode: VideoRenderMode.RENDER_MODE_HIDDEN,
          });
        } catch (error) {
          console.log('error', error);
        }
      });
  }

  changeAudioState(isMicOn: boolean): void {
    isMicOn ? this.engine.startAudioCapture() : this.engine.stopAudioCapture();
  }

  async changeVideoState(isVideoOn: boolean): Promise<void> {
    isVideoOn
      ? await this.engine.startVideoCapture()
      : await this.engine.stopVideoCapture();
    this.engine.setLocalVideoMirrorType(1);
  }

  leave(): void {
    this.engine.leaveRoom();
    VERTC.destroyEngine(this.engine);
  }

  destroy(): void {
    this.engine.stopVideoCapture();
    this.engine.stopAudioCapture();
  }

  createEngine(): void {
    this.engine = VERTC.createEngine(this.config.appId);
  }

  unpublish(): Promise<void> {
    return this.engine.unpublish();
  }

  /**
   * @brief 共享屏幕音视频流
   * @function createScreenStream
   * @param {*} screenConfig
   */
  async createScreenStream(
    screenStreamSettings = {},
    success: () => void,
    fail: (err: Error) => void
  ): Promise<void> {
    try {
      await this.engine.startScreenCapture();
      this.engine
        .publishScreen()
        .then(() => {
          success();
        })
        .catch((err: Error) => console.log('err', err));
    } catch (e: any) {
      fail(e);
    }
  }
  /**
   * @brief 取消共享屏幕音视频流
   * @function stopScreenStream
   */
  async stopScreenStream(success: () => void): Promise<void> {
    await this.engine.stopScreenCapture();
    await this.engine.unpublishScreen();
    success();
  }
}
