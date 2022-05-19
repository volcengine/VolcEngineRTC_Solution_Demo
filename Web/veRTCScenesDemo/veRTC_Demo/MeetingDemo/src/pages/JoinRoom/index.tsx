import React, { Component, ReactNode } from 'react';
import { history } from 'umi';
import { FormInstance } from 'antd/es/form';
import { Form, Input, Button, Tooltip, Modal } from 'antd';
import { WrappedComponentProps } from 'react-intl';
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
import { LocalPlayer } from '../Meeting/components/MediaPlayer';
import { injectProps, ConnectedProps, connector } from '../Meeting/configs/config';
import DeviceController from '@/lib/DeviceController';

const logger = new Logger('JoinRoom');

export type LoginProps = ConnectedProps<typeof connector> &
  WrappedComponentProps;

export interface LoginState {
  nameErrType: ERROR_TYPES;
  roomIdErrType: ERROR_TYPES;
  settingsVisible: boolean;
  cameraStream: {
    playerComp: JSX.Element;
  } | null;
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

  /**
   *设备通用方法
   */
  deviceLib = new DeviceController(this.props);

  get nameHasErr(): boolean {
    return this.state.nameErrType !== ERROR_TYPES.VALID;
  }

  get roomIdHasErr(): boolean {
    return this.state.roomIdErrType !== ERROR_TYPES.VALID;
  }

  componentDidMount(): void {
    const {
      currentUser,
      settings
    } = this.props;

    const param = {
      currentUser,
      settings,
    };
    this.props.rtc.createEngine();
    this.deviceLib.openCamera(param, () => {
      this.props.setLocalAudioVideoCaptureSuccess(true);
      this.setState({
        cameraStream: {
          playerComp: (
            <LocalPlayer
              localCaptureSuccess={true}
              rtc={this.props.rtc}
              renderDom="local-preview-player"
            />
          ),
        },
      });
    }, false);
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
            this.setState({ cameraStream: null });
          }
          const loginInfo = Utils.getLoginInfo();
          if (loginInfo) {
            loginInfo.user_name = userName;
            Utils.setLoginInfo(loginInfo);
          }
          this.props.setMeetingStatus('start');
          history.push(`/meeting?roomId=${roomId}&username=${userName}`);
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
      isMicOn,
    } = this.props.currentUser;
    if (!audio) {
      audioMessage && modalWarning(audioMessage);
      return;
    }
    this.props.rtc.changeAudioState(!isMicOn);
    this.props.setIsMicOn(!isMicOn);
  };

  onVideoClick = (): void => {
    const {
      deviceAccess: { video, videoMessage },
      isCameraOn,
    } = this.props.currentUser;
    if (!video) {
      videoMessage && modalWarning(videoMessage);
      return;
    }
    this.props.rtc.changeVideoState(!isCameraOn);
    this.props.setIsCameraOn(!isCameraOn);
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
      currentUser: { isCameraOn, isMicOn },
      meeting: { status },
    } = this.props;

    return (
      <div className={styles['join-room-container']}>
        <View
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
          volume={0}
          sharingView={false}
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
                    disabled={
                      !this.props.currentUser.network || !roomId || !userName
                    }
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
          {process.env.NODE_ENV !== 'development' && (
            <FeedBack status={status} />
          )}
        </div>
        {process.env.NODE_ENV !== 'development' && (
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
        )}
      </div>
    );
  }
}

export default injectProps(Login);
