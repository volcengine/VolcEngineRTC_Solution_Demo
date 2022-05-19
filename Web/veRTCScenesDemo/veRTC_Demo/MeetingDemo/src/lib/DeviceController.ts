import { StreamIndex } from '@volcengine/rtc';
import {
  MeetingProps,
} from '@/app-interfaces';
import { TOASTS } from '@/constant';
import { message } from 'antd';
import { UserModelState } from '../models/user';
import { MeetingSettingsState } from '../models/meeting-settings';
import { ICreateStreamRes } from '../rtcApi/types';
import { MeetingModelState } from '../models/meeting';


/**
 * @brief 设备管理
 * @function DeviceController
 */
class DeviceController {
  private props: MeetingProps;

  constructor(props: MeetingProps) {
    this.props = props;
  }

  /**
   * @brief 设置属性
   * @function setProps
   */
  public setProps = (props: MeetingProps): void => {
    this.props = props;
  };

  /**
   * @brief 打开或关闭音频设备
   * @function changeAudioState
   */
  public changeAudioState = (state: boolean): void => {
    state ? this.props.mc?.turnOnMic({}) : this.props.mc?.turnOffMic({});
    this.props.rtc.changeAudioState(state);
    this.props.setIsMicOn(state);
  };

  /**
   * @brief 打开或关闭摄像头设备
   * @function changeVideoState
   */
  public changeVideoState = (state: boolean): void => {
    state ? this.props.mc?.turnOnCamera({}) : this.props.mc?.turnOffCamera({});
    this.props.rtc.changeVideoState(state);
    this.props.setIsCameraOn(state);
  };

  public openCamera(
    param: {
      currentUser: UserModelState;
      settings: MeetingSettingsState;
    },
    callback: () => void,
    isPublish: boolean
  ): void {
    const { currentUser, settings } = param;

    const { isCameraOn, isMicOn } = currentUser;
    const streamOptions = {
      mic: settings.mic,
      camera: settings.camera,
      video: isCameraOn,
      audio: isMicOn,
    };
    const { engine } = this.props.rtc;
    const { streamSettings, screenStreamSettings } = settings;

    engine.setVideoCaptureConfig({
      frameRate: streamSettings.frameRate.max,
      ...streamSettings.resolution,
    });
    engine.setVideoEncoderConfig(StreamIndex.STREAM_INDEX_MAIN, [
      {
        maxKbps: streamSettings.bitrate.max,
      },
    ]);
    engine.setVideoEncoderConfig(StreamIndex.STREAM_INDEX_SCREEN, [
      {
        maxKbps: screenStreamSettings.bitrate.max,
      },
    ]);
    //TODO：选择设备、获取本地预览流
    this.props.rtc?.createLocalStream(
      streamOptions,
      isPublish,
      (res: ICreateStreamRes) => {
        const { code, msg, devicesStatus } = res;
        if (code === -1) {
          alert(msg);
          this.props.setIsMicOn(false);
          this.props.setIsCameraOn(false);
        } else {
          let param = currentUser.deviceAccess;
          if (devicesStatus['audio'] === 0) {
            this.props.setIsMicOn(false);
            param = {
              ...param,
              audio: false,
              audioMessage: 'mic_setting_right',
            };
          }
          if (devicesStatus['video'] === 0) {
            this.props.setIsCameraOn(false);
            param = {
              ...param,
              video: false,
              videoMessage: 'car_setting_right',
            };
          }
          this.props.setDeviceAccess({
            ...param,
          });
        }
        callback();
      }
    );
  }
  /**
   * @brief 屏幕分享
   * @function changeShareState
   */
  public changeShareState(
    param: {
      meeting: MeetingModelState;
      settings: MeetingSettingsState;
    },
    isShare: boolean
  ): void {
    const { meeting, settings } = param;

    if (isShare && meeting.isSharing && meeting.meetingInfo.screen_shared_uid) {
      message.warn('当前有人正在分享，请等待结束后再开启');
      return;
    }
    if (isShare) {
      // 开始共享
      const { frameRate, bitrate, resolution } = settings.screenStreamSettings;

      this.props.rtc.createScreenStream(
        {
          width: resolution.width,
          height: resolution.height,
          frameRate: frameRate.max,
          maxKbps: bitrate.max,
        },
        () => {
          this.props.setIsSharing(isShare);
          this.props.setMeetingIsSharing(true);
          this.props.mc?.startShareScreen({});
        },
        (err: {
          error: { code: number; message: string; name: string };
          message: string;
          code: string;
        }) => {
          if (this.props.currentUser.isSharing) {
            this.stopShare(false);
          }
          if (err.error.name === 'NotAllowedError')
            message.error(TOASTS['screen_not_allow']);
          else message.error(TOASTS['screen_error']);
        }
      );
    } else {
      // 结束共享
      this.stopShare(isShare);
    }
  }
  public stopShare(isShare: boolean): void {
    this.props.setMeetingIsSharing(false);
    this.props.mc?.endShareScreen();
    this.props.rtc.destoryScreenStream(() => {
      this.props.setIsSharing(isShare);
    });
  }
}

export default DeviceController;
