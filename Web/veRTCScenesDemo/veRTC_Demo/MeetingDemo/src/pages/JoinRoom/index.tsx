import React, { Component, ReactNode } from 'react';
import RTC from '@/sdk/VRTC.esm.min.js';
import { Dispatch } from '@@/plugin-dva/connect';
import { injectIntl, history } from 'umi';
import { connect, bindActionCreators } from 'dva';
import { ConnectedProps } from 'react-redux';
import { FormInstance } from 'antd/es/form';
import { Form, Input, Button, Tooltip, Modal } from 'antd';
import { WrappedComponentProps } from 'react-intl';
import { userActions } from '@/models/user';
import { meetingActions } from '@/models/meeting';
import { AppState, Stream } from '@/app-interfaces';
import Logger from '@/utils/Logger';
import SettingsModal from '@/components/SettingsModal';
import FeedBack from '@/components/FeedBack';
import styles from './index.less';
import IconBtn from '@/components/IconBtn';
import View from '@/components/View';
import micOnIon from '/assets/images/micOnIon.png';
import micOffIcon from '/assets/images/micOffIcon.png';
import camOnIcon from '/assets/images/camOnIcon.png';
import camOffIcon from '/assets/images/camOffIcon.png';
import settingsIcon from '/assets/images/settingsIcon.png';
import camPause from '/assets/images/camPause.png';
import icon from '/assets/images/icon-256@2x.png';
import Utils from '@/utils/utils';
import { StoreValue } from 'rc-field-form/lib/interface';
import { modalWarning } from '@/pages/Meeting/components/MessageTips';
import Player from '../Meeting/components/Player';

function mapStateToProps(state: AppState) {
  return {
    user: state.user,
    settings: state.meetingSettings,
    meeting: state.meeting,
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators({ ...userActions, ...meetingActions }, dispatch),
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

const logger = new Logger('JoinRoom');

export type LoginProps = ConnectedProps<typeof connector> &
  WrappedComponentProps;

export interface LoginState {
  nameErrType: ERROR_TYPES;
  roomIdErrType: ERROR_TYPES;
  settingsVisible: boolean;
  cameraStream: Stream | null;
}

enum ERROR_TYPES {
  VALID,
  EMPTY_STRING,
  INVALID_CHARACTERS,
}

export class Login extends Component<LoginProps, LoginState> {
  constructor(props: LoginProps) {
    super(props);
    this.state = {
      nameErrType: ERROR_TYPES.VALID,
      roomIdErrType: ERROR_TYPES.VALID,
      settingsVisible: false,
      cameraStream: null,
    };
  }

  formRef = React.createRef<FormInstance>();

  get nameHasErr(): boolean {
    return this.state.nameErrType !== ERROR_TYPES.VALID;
  }

  get roomIdHasErr(): boolean {
    return this.state.roomIdErrType !== ERROR_TYPES.VALID;
  }

  componentDidUpdate(prevProps: LoginProps): void {
    if ( this.props.settings && prevProps.settings.streamSettings !== this.props.settings.streamSettings) {
      if (this.state.cameraStream) {
        this.state.cameraStream.setVideoEncoderConfiguration(
          this.props.settings.streamSettings
        );
      }
    }

    if (this.props.settings?.camera !== prevProps.settings?.camera || this.props.settings?.mic !== prevProps.settings?.mic) {
      this.openCamera();
    }
  }

  openCamera(): void{
    const {
      user: { isCameraOn, isMicOn },
      settings,
    } = this.props;

    const stream = RTC.createStream({
      video: true,
      audio: true,
      microphoneId: settings?.mic,
      cameraId: settings?.camera,
    });

    stream.setVideoEncoderConfiguration(settings.streamSettings);

    stream.init(() => {
      logger.debug('stream init success');
      if (isCameraOn) {
        stream.unmuteVideo();
      } else {
        stream.muteVideo();
      }
      if (isMicOn) {
        stream.unmuteAudio();
      } else {
        stream.muteAudio();
      }

      const player = <Player stream={stream} />;
      stream.playerComp = player;

      this.setState({ cameraStream: stream });
    });
  }

  componentDidMount(): void {

    this.openCamera();

    RTC.checkAudioPermission(
      () => {
        this.props.setDeviceAccess({
          ...this.props.user.deviceAccess,
          audio: true,
        });
      },
      (err: Error) => {
        this.props.setIsMicOn(false);
        this.props.setDeviceAccess({
          ...this.props.user.deviceAccess,
          audio: false,
          audioMessage: err.toString()?.includes('Permission')
            ? 'mic_right'
            : 'mic_setting_right',
        });
      }
    );

    RTC.checkVideoPermission(
      () => {
        this.props.setDeviceAccess({
          ...this.props.user.deviceAccess,
          video: true,
        });
      },
      (err: Error) => {
        this.props.setIsCameraOn(false);
        this.props.setDeviceAccess({
          ...this.props.user.deviceAccess,
          video: false,
          videoMessage: err.toString()?.includes('Permission denied')
            ? 'car_right'
            : 'car_setting_right',
        });
      }
    );
  }

  /**
   * Rules for login fields validation
   * @param {"name" | "roomId"} name
   * @return {Rule[]}
   */
  getLoginFieldRules = (
    value: StoreValue,
    name: 'name' | 'roomId',
    regRes: boolean
  ): Promise<void | any> | void => {
    const errorTypeKey = name === 'name' ? 'nameErrType' : 'roomIdErrType';

    const setStateObj = {
      [errorTypeKey]: ERROR_TYPES.VALID,
    } as {
      [K in 'nameErrType' | 'roomIdErrType']: ERROR_TYPES;
    };

    let result: Promise<Error | void>;

    if (!value || regRes) {
      setStateObj[errorTypeKey] = !value
        ? ERROR_TYPES.EMPTY_STRING
        : ERROR_TYPES.INVALID_CHARACTERS;
      result = Promise.reject(new Error(' '));
    } else {
      result = Promise.resolve();
    }

    this.setState(setStateObj);

    return result;
  };

  onLogin = async (): Promise<void> => {
    if (this.formRef.current) {
      try {
        this.formRef.current.validateFields().then((values) => {
          const { userName, roomId } = values;
          this.props.setRoomId(roomId);
          this.props.setUserName(userName);
          if (this.state.cameraStream) {
            this.state.cameraStream.close();
            this.setState({ cameraStream: null });
          }
          const loginInfo = Utils.getLoginInfo();
          if (loginInfo) {
            loginInfo.user_name = userName;
            Utils.setLoginInfo(loginInfo);
          }
          this.props.setMeetingStatus('start');
          history.push(
            `/meeting?roomId=${roomId}&username=${userName}`
          );
        });
      } catch (e) {
        logger.error(e);
        // TODO
      }
    }
  };

  onAudioClick = (): void => {
    const {
      deviceAccess: { audio, audioMessage },
    } = this.props.user;

    if (!audio) {
      audioMessage && modalWarning(audioMessage);
      return;
    }
    this.props.setIsMicOn(!this.props.user.isMicOn);
  };

  onVideoClick = (): void => {
    const {
      deviceAccess: { video, videoMessage },
    } = this.props.user;
    if (!video) {
      videoMessage && modalWarning(videoMessage);
      return;
    }
    this.props.setIsCameraOn(!this.props.user.isCameraOn);
  };

  openSettings = (): void => {
    this.setState({
      settingsVisible: true,
    });
  };

  renderToolTip = (): ReactNode => {
    return (
      <div>
        <strong>非法输入，输入规则如下：</strong>
        <div>1. 26个大写字母 A ~ Z 。</div>
        <div>2. 26个小写字母 a ~ z 。</div>
        <div>3. 10个数字 0 ~ 9 。</div>
        <div>
          4. 下划线&quot;_&quot;, at符&quot;@&quot;, 减号&quot;-&quot;。
        </div>
      </div>
    );
  };

  render(): ReactNode {
    const {
      user: { isCameraOn, isMicOn },
      meeting: { status },
    } = this.props;

    return (
      <div className={styles['join-room-container']}>
        <View
          stream={this.state.cameraStream}
          player={this.state.cameraStream?.playerComp}
          is_host={false}
          is_sharing={false}
          is_camera_on={isCameraOn}
          is_mic_on={isMicOn}
          created_at={Date.now()}
          room_id=""
          user_id="void"
          user_name=""
          user_uniform_id=""
          avatarOnCamOff={<img style={{ width: '50%' }} src={camPause} />}
        />

        <div className={styles['login-toolbar']}>
          <Form
            layout="inline"
            ref={this.formRef}
            style={{
              height: 32,
            }}
          >
            <div className={styles['login-toolbar-title']}>登录</div>
            <Tooltip
              visible={this.roomIdHasErr}
              title={
                this.state.roomIdErrType === ERROR_TYPES.INVALID_CHARACTERS
                  ? this.renderToolTip()
                  : '请填写房间ID'
              }
              placement="topLeft"
            >
              <Form.Item
                name="roomId"
                validateTrigger="onChange"
                initialValue={Utils.getQueryString('roomId')}
                rules={[
                  {
                    validator: (_, value) => {
                      const regRes = !/^[0-9a-zA-Z_\-@.]*$/.test(value);
                      return this.getLoginFieldRules(value, 'roomId', regRes);
                    },
                  },
                ]}
              >
                <Input
                  placeholder={'房间ID'}
                  maxLength={18}
                  style={{
                    width: 160,
                  }}
                />
              </Form.Item>
            </Tooltip>
            <Tooltip
              visible={this.nameHasErr}
              title={
                this.state.nameErrType === ERROR_TYPES.INVALID_CHARACTERS
                  ? this.renderToolTip()
                  : '请填写用户名'
              }
              placement="topLeft"
            >
              <Form.Item
                initialValue={Utils.getLoginUserName()}
                name="userName"
                rules={[
                  {
                    validator: (_, value) => {
                      const regRes = !/^[0-9a-zA-Z_\-@.\u4e00-\u9fa5]*$/.test(
                        value
                      );
                      return this.getLoginFieldRules(value, 'name', regRes);
                    },
                  },
                ]}
              >
                <Input
                  placeholder={'用户名'}
                  maxLength={18}
                  style={{
                    width: 160,
                  }}
                />
              </Form.Item>
            </Tooltip>
            <Form.Item
              shouldUpdate={(prevValues, curValues) =>
                prevValues.userName !== curValues.userName ||
                prevValues.roomId !== curValues.roomId
              }
            >
              {({ getFieldsValue }) => {
                const { roomId, userName } = getFieldsValue();
                return (
                  <Button
                    type={'primary'}
                    onClick={this.onLogin}
                    disabled={!this.props.user.network || !roomId || !userName}
                  >
                    进入房间
                  </Button>
                );
              }}
            </Form.Item>
          </Form>
          <IconBtn
            width={32}
            height={32}
            onClick={this.onAudioClick}
            style={{ marginLeft: 32 }}
          >
            <img
              src={isMicOn ? micOnIon : micOffIcon}
              alt={`mic-${isMicOn ? 'on' : 'off'}`}
            />
          </IconBtn>
          <IconBtn width={32} height={32} onClick={this.onVideoClick}>
            <img
              src={isCameraOn ? camOnIcon : camOffIcon}
              style={{
                width: isCameraOn ? 18 : undefined,
                height: isCameraOn ? 12 : undefined,
              }}
              alt={`camera-${isCameraOn ? 'on' : 'off'}`}
            />
          </IconBtn>
          <IconBtn width={32} height={32} onClick={this.openSettings}>
            <img src={settingsIcon} alt="settings" />
          </IconBtn>
          {this.state.settingsVisible ? (
            <SettingsModal
              visible={this.state.settingsVisible}
              close={() => this.setState({ settingsVisible: false })}
            />
          ) : undefined}
          <FeedBack status={status} />
        </div>
        <Modal
          visible={status === 'init'}
          footer={null}
          closable={false}
          maskClosable={false}
          width={280}
          className={styles['modal']}
          bodyStyle={{ fontSize: 12, background: '#ddd' }}
        >
          <div className={styles['icon']}>
            <img src={icon} />
          </div>
          <p className={styles['tips']}>
            本产品仅用于功能体验，单次会议时长不超过10分钟
          </p>
          <Button
            onClick={() => this.props.setMeetingStatus('closeTips')}
            type="primary"
            block
          >
            确定
          </Button>
        </Modal>
      </div>
    );
  }
}

export default connector(injectIntl(Login));
