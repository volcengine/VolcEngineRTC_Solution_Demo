import React, { Component, ReactNode } from 'react';
import RTC from '@/sdk/VRTC.esm.min.js';
import { message, Skeleton } from 'antd';
import { injectIntl, history } from 'umi';
import { connect, bindActionCreators } from 'dva';
import { ConnectedProps } from 'react-redux';
import { throttle, concat } from 'lodash';
import { WrappedComponentProps } from 'react-intl';
import { Dispatch } from '@@/plugin-dva/connect';
import Logger from '@/utils/Logger';
import Utils from '@/utils/utils';
import { userActions } from '@/models/user';
import { meetingActions, ViewMode } from '@/models/meeting';
import Header from './components/Header';
import MeetingViews from './components/MeetingViews';
import ControlBar from './components/ControlBar';
import UsersDrawer from './components/UsersDrawer';
import LeavingConfirm from './components/LeavingConfirm';
import StreamStats from './components/StreamStats';
import { modalError } from './components/MessageTips';
import Player from './components/Player';
import {
  showOpenMicConfirm,
  sendMutedInfo,
  sendInfo,
  hostChangeInfo,
} from './components/MessageTips';

import { AppState, Stream, IRemoteAudioLevel } from '@/app-interfaces';
import type { MeetingUser } from '@/models/meeting';
import type {
  GetAppIDResponse,
  JoinMeetingResponse,
  UserPayload,
  UserStatusChangePayload,
  HostChangePayload,
} from '@/lib/socket-interfaces';
import styles from './index.less';
import { TOASTS } from '@/config';

const logger = new Logger('meeting');

function mapStateToProps(state: AppState) {
  return {
    currentUser: state.user,
    meeting: state.meeting,
    mc: state.meetingControl.sdk,
    settings: state.meetingSettings,
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators({ ...userActions, ...meetingActions }, dispatch),
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export type IMeetingProps = ConnectedProps<typeof connector> &
  WrappedComponentProps;

export interface IMeetingState {
  usersDrawerVisible: boolean;
  cameraStream: Stream | null;
  screenStream: Stream | null;
  remoteStreams: { [id: string]: Stream };
  leaving: boolean;
  audioLevels: IRemoteAudioLevel[];
  localVolume: number;
  refresh: boolean;
}

const initState = {
  usersDrawerVisible: false,
  cameraStream: null,
  screenStream: null,
  remoteStreams: {},
  leaving: false,
  audioLevels: [],
  localVolume: 0,
  refresh: false,
};

class Meeting extends Component<IMeetingProps, IMeetingState> {
  constructor(props: IMeetingProps) {
    super(props);

    window.__meeting = this;

    this.state = initState;

    this.hanldleMeetingController = this.hanldleMeetingController.bind(this);
    this.unMount = this.unMount.bind(this);
    this.hanldleMeetingController();
  }

  intervalId: undefined| ReturnType<typeof setTimeout>;

  componentDidMount(): void {
    const debug = Utils.getQueryString('debug');
    if (debug) {
      localStorage.setItem('debug', 'veRTCDemo:*');
    } else {
      localStorage.removeItem('debug');
    }
    logger.debug('meeting props: %o', this.props);

    if (!this.props.mc) {
      logger.warn('joinMeeting before meeting control init !');
      return;
    }

    this.props.mc.checkSocket().then(() => {
      if (this.state.refresh) {
        setTimeout(() => {
          this.joinMeeting();
          this.setState({
            refresh: false,
          });
        }, 500);
      } else {
        this.joinMeeting();
      }
      this.handleMeetingStatus();
    });

    window.addEventListener('beforeunload', this.unMount);

    this.checkMediaState();

    history.listen((location, action) => {
      if (action == 'POP') {
        if (location.pathname === '/meeting') {
          this.setState({
            refresh: true,
          });
        }
      }
    });
  }

  componentDidUpdate(prevProps: IMeetingProps, preState: IMeetingState): void {
    if (
      this.props.mc &&
      prevProps.currentUser.isHost !== this.props.currentUser.isHost
    ) {
      this.props.mc.isHost = this.props.currentUser.isHost;
    }
    if (
      this.props.settings &&
      prevProps.settings.streamSettings !== this.props.settings.streamSettings
    ) {
      if (this.state.cameraStream) {
        this.state.cameraStream.setVideoEncoderConfiguration(
          this.props.settings.streamSettings
        );
      }
    }
    if (
      this.props.settings &&
      prevProps.settings.screenStreamSettings !==
        this.props.settings.screenStreamSettings
    ) {
      if (this.state.screenStream) {
        this.state.screenStream.setVideoEncoderConfiguration(
          this.props.settings.screenStreamSettings
        );
      }
    }
    if (
      this.props.settings?.camera !== prevProps.settings?.camera ||
      this.props.settings?.mic !== prevProps.settings?.mic
    ) {
      const { cameraStream } = this.state;
      if (cameraStream) {
        this.props.mc?.unpublish(cameraStream as Stream).then(() => {
          this.state.cameraStream?.close();
          this.openCamera({
            video: this.props.currentUser.isCameraOn,
            audio: this.props.currentUser.isMicOn,
          });
        });
      }
    }
    if (this.state.audioLevels !== preState?.audioLevels) {
      this.props.meeting.status !== 'end' && this.debouncedChangeMeetingUserOrder();
    }
  }

  unMount() {
    this.props.setMeetingUsers([]);
    this.props?.mc?.removeEvent();
    this.closeAllStream();
    this.props?.mc?.cleanStreamRecord();
    this.props.setMeetingUsers([]);
    this.props.setMeetingOrderUsers([]);
    this.props.setViewMode(ViewMode.GalleryView);
    this.props.setSpeakCollapse(false);
    this.setState({
      ...initState,
    });
    this.intervalId && clearInterval(this.intervalId);
  }

  componentWillUnmount() {
    this.unMount();
    window.removeEventListener('beforeunload', this.unMount);
  }

  handleMeetingStatus() {
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.props.setMeetingStatus('hidden');
      }
    });
  }

  hanldleMeetingController() {
    const { setMeetingUsers, setIsHost, setViewMode, setMeetingInfo } =
      this.props;

    // broadcasting / 广播消息
    this.props.mc?.on(
      'onUserMicStatusChange',
      (payload: UserStatusChangePayload) => {
        this.changeUserStatus('is_mic_on', payload);
        setTimeout(() => {
          this.getRemoteStreamStat(true);
        }, 1000);
      }
    );

    this.props.mc?.on(
      'onUserCameraStatusChange',
      (payload: UserStatusChangePayload) => {
        this.changeUserStatus('is_camera_on', payload);
      }
    );

    this.props.mc?.on('onHostChange', (payload: HostChangePayload) => {
      const users = this.props.meeting.meetingUsers.map((user) => {
        const ret = { ...user };
        if (ret.user_id === payload.host_id) {
          ret.is_host = true;
        } else {
          ret.is_host = false;
        }
        return ret;
      });
      setMeetingUsers(users);
      const isHost = payload.host_id === this.props.currentUser.userId;
      if (isHost) {
        message.info('你已成为主持人');
      }
      setIsHost(isHost);
    });

    this.props.mc?.on('onUserJoinMeeting', (payload: MeetingUser) => {
      const users = this.props.meeting.meetingUsers;
      if (users.some((user) => user.user_id === payload.user_id)) {
        logger.error('onUserJoinMeeting 该用户 %s 已存在', payload.user_id);
      } else {
        setMeetingUsers([...users, payload]);
      }
    });

    this.props.mc?.on('onUserLeaveMeeting', (payload: { user_id: string }) => {
      const users = this.props.meeting.meetingUsers.filter(
        (user) => user.user_id !== payload.user_id
      );
      setMeetingUsers([...users]);
    });

    this.props.mc?.on(
      'onShareScreenStatusChanged',
      (payload: UserStatusChangePayload) => {
        if (payload.status) {
          setViewMode(ViewMode.SpeakerView);
          setMeetingInfo({
            ...this.props.meeting.meetingInfo,
            screen_shared_uid: payload.user_id,
          });
        } else {
          setViewMode(ViewMode.GalleryView);
          setMeetingInfo({
            ...this.props.meeting.meetingInfo,
            screen_shared_uid: '',
          });
        }
        this.changeUserStatus('is_sharing', payload);
      }
    );

    this.props.mc?.on('onRecord', () => {
      setMeetingInfo({
        ...this.props.meeting.meetingInfo,
        record: true,
      });
    });

    this.props.mc?.on('onMeetingEnd', () => {
      logger.debug('onMeetingEnd');
      this.props.setMeetingStatus('end');
      this.props.mc?.unpublish(this.state.cameraStream as Stream).then(() => {
        this.props.mc?.clientLeave(() => this.end());
      });
    });

    this.props.mc?.on('onMuteAll', () => {
      if (this.state.cameraStream && !this.props.currentUser.isHost) {
        this.state.cameraStream.disableAudio();
        this.props.setIsMicOn(false);
      }
      //user who is not host will receive muted message
      if (!this.props.currentUser.isHost) {
        this.props.setIsMicOn(false);
        sendMutedInfo();
      }
      if (this.props.currentUser?.userId) {
        this.props.mc
          ?.getMeetingUserInfo({
            user_id: this.props.currentUser?.userId,
          })
          .then((res) => {
            this.props.setMeetingUsers(res);
          });
      }
    });

    // user singel message
    this.props.mc?.on('onMuteUser', (payload: UserPayload) => {
      sendMutedInfo();
      this.changeMicState(false);
    });

    this.props.mc?.on('onAskingMicOn', (payload: UserPayload) => {
      showOpenMicConfirm(() => this.changeMicState(true));
    });

    this.props.mc?.on('onUserKickedOff', () => {
      message.warning(TOASTS['tick']);
      this.end();
    });

    // rtc events
    this.props.mc?.on('OnReceivedStream', (stream: Stream) => {
      logger.debug('OnReceivedStream');
      const player = <Player stream={stream} remote={true} />;
      stream.playerComp = player;

      if (stream.stream.screen) {
        this.setState({
          screenStream: stream,
        });
      } else {
        const _s = {
          ...this.state.remoteStreams,
          [stream.uid]: stream,
        };

        this.setState({
          remoteStreams: _s,
        });

        if (!this.props.mc) return;
        this.props.mc.streams = _s;

        this.getRemoteStreamStat(true);
      }
    });

    this.props.mc?.on(
      'onRemoveStream',
      ({ uid, screen }: { uid: string; screen: boolean }) => {
        logger.debug('onRemoveStream');

        if (screen) {
          this.state.screenStream?.close();
          this.setState({ screenStream: null });
        } else {
          const _s = this.state.remoteStreams;

          Object.keys(_s).map((key) => {
            const stream = _s[key];
            if (stream.getId() === uid) {
              stream?.close();
              delete _s[key];
            }
          });

          this.setState({
            remoteStreams: _s,
          });

          if (!this.props.mc) return;

          this.props.mc.streams = _s;
        }
      }
    );

    this.props.mc?.on('onUserReconnect', (code: number) => {
      if (code === 200) {
        if (this.props.currentUser?.userId) {
          this.props.mc
            ?.getMeetingInfo()
            .then(() => {
              this.props.mc
                ?.getMeetingUserInfo({})
                .then((res) => {
                  this.props.setMeetingUsers(res);
                })
                .catch((err) => {
                  if (err === 'record not found') {
                    message.warning('会议已结束', undefined, () => {
                      this.end();
                    });
                  }
                });
            })
            .catch((err) => {
              if (err === 'record not found') {
                message.warning('会议已结束', undefined, () => {
                  this.end();
                });
              }
            });
        }
      } else {
        if (code === 404) {
          message.warning('你已被踢出房间', undefined, () => {
            this.end();
          });
          return;
        }
        if (code === 422) {
          message.warning('会议已结束', undefined, () => {
            this.end();
          });
          return;
        }
      }
    });
  }

  debouncedChangeMeetingUserOrder = throttle(
    () => this.changeMeetingUserOrder(),
    1000 * 8
  );

  changeMeetingUserOrder() {
    const { meeting, currentUser } = this.props;
    const { audioLevels } = this.state;

    if (meeting.meetingUsers?.length) {
      logger.debug('audioLevels %', audioLevels);

      const cloneMeeting = [...meeting?.meetingUsers];

      let hostAndMe = [];

      const hostIndex = cloneMeeting.findIndex((item) => item.is_host);
      const host = cloneMeeting.splice(hostIndex, 1);

      if (currentUser?.isHost) {
        hostAndMe = host;
      } else {
        const meIndex = cloneMeeting.findIndex(
          (item) => item.user_id === currentUser.userId
        );
        const me = cloneMeeting.splice(meIndex, 1);
        hostAndMe = concat([], host, me);
      }

      if (audioLevels.length === 0) {
        this.props.setMeetingOrderUsers(concat([], hostAndMe, cloneMeeting));
        return;
      }

      const audioLevels_userId: string[] = audioLevels.map(
        (item) => item?.userId || ''
      );

      cloneMeeting.sort((a, b) => {
        if (audioLevels_userId.indexOf(a?.user_id) === -1) {
          return 1;
        }
        if (audioLevels_userId.indexOf(b?.user_id) === -1) {
          return -1;
        }
        return (
          audioLevels_userId.indexOf(a?.user_id) -
          audioLevels_userId.indexOf(b?.user_id)
        );
      });

      const _combine = concat([], hostAndMe, cloneMeeting);

      this.props.setMeetingOrderUsers(_combine);

      this.props.mc?.changSubscribe(_combine);
    }
  }

  async getRemoteStreamStat(forceUpdate?: boolean): Promise<void> {
    if (!this.props.mc) return;

    const { remoteStreams, cameraStream } = this.state;

    logger.debug('remoteStreams %', remoteStreams);

    const remotesAudioStats = await this.props.mc?.getRemoteAudioStats();
    const volume = (cameraStream?.getAudioLevel() ?? 0) / 10;

    if (remotesAudioStats) {
      const remoteSort = Object.entries(remotesAudioStats)
        ?.map(([key, value]: [string, { RecvLevel: number }]) => ({
          userId: key,
          RecvLevel: value.RecvLevel,
        }))
        .sort((a, b) => (b.RecvLevel || 0) - (a.RecvLevel || 0));

      this.setState({
        audioLevels: remoteSort,
        localVolume: volume,
      });

      if (forceUpdate) {
        this.changeMeetingUserOrder();
      }
    }
  }

  joinMeeting(): void {
    const {
      currentUser,
      setMeetingInfo,
      setMeetingUsers,
      setIsHost,
      setViewMode,
    } = this.props;

    const userId = Utils.getLoginUserId();
    this.props.setUserId(userId);

    this.props.mc
      ?.getAppID({})
      .then((app?: GetAppIDResponse) => {
        if (!app) {
          return;
        }
        return this.props.mc?.joinMeeting({
          app_id: app.app_id,
          user_id: userId,
          user_name: Utils.getLoginUserName(),
          room_id: this.roomId,
          mic: currentUser.isMicOn,
          camera: currentUser.isCameraOn,
        });
      })
      .then((res?: JoinMeetingResponse) => {
        if (!res) {
          return;
        }
        const { info, users } = res;

        setMeetingInfo(info);
        setMeetingUsers(users);
        if (info.screen_shared_uid) {
          setViewMode(ViewMode.SpeakerView);
        }
        setIsHost(info.host_id === Utils.getLoginUserId());

        if (!this.state.cameraStream) {
          this.openCamera(
            {
              video: this.props.currentUser.isCameraOn,
              audio: this.props.currentUser.isMicOn,
            },
            () => {
              this.intervalId = setInterval(() => {
                this.getRemoteStreamStat();
              }, 1000 * 2);
              this.changeMeetingUserOrder();
            }
          );
        }
      })
      .catch((e) => {
        message.error(`Join meeting failed: ${e || 'unknown'}`, () => {
          if (e === 'login token expired') {
            Utils.removeLoginInfo();
            history.push('/login');
          } else {
            this.end();
          }
        });
      });
  }

  changeMicState(s: boolean): void {
    logger.debug('changeMicState: ', s);
    const { setIsMicOn } = this.props;
    const { cameraStream } = this.state;
    setIsMicOn(s);
    s
      ? this.props.mc?.turnOnMic({})
      : this.props.mc
          ?.turnOffMic({})
          .catch(() => message.error(TOASTS['mute_error']));
    if (cameraStream) {
      if (s) {
        cameraStream.enableAudio();
      } else {
        cameraStream.disableAudio();
      }
    } else {
      this.openCamera({
        audio: this.props.currentUser?.isMicOn,
      });
    }
  }

  changeUserStatus(
    state: 'is_mic_on' | 'is_camera_on' | 'is_sharing',
    payload: UserStatusChangePayload
  ): void {
    logger.debug(payload);
    const users = this.props.meeting.meetingUsers.map((user) => {
      const ret = { ...user };
      if (ret.user_id === payload.user_id) {
        ret[state] = payload.status;
      }
      return ret;
    });
    this.props.setMeetingUsers(users);
  }

  changeCameraState(s: boolean): void {
    logger.debug('changeCameraState: ', s);
    const { setIsCameraOn } = this.props;
    const { cameraStream } = this.state;
    setIsCameraOn(s);
    s ? this.props.mc?.turnOnCamera({}) : this.props.mc?.turnOffCamera({});
    if (cameraStream) {
      if (s) {
        cameraStream.enableVideo();
      } else {
        cameraStream.disableVideo();
      }
    } else {
      this.openCamera({
        video: this.props.currentUser.isCameraOn,
      });
    }
  }

  openCamera(
    constraints: { video?: boolean; audio?: boolean },
    callback?: () => void
  ): void {
    const { settings, currentUser } = this.props;

    const stream = RTC.createStream({
      video: currentUser?.deviceAccess?.video,
      audio: currentUser?.deviceAccess?.audio,
      microphoneId: settings?.mic,
      cameraId: settings?.camera,
    });

    stream.setVideoEncoderConfiguration(settings.streamSettings);

    stream.init(
      () => {
        logger.debug('stream init success');
        if (constraints.video) {
          stream.unmuteVideo();
        } else {
          stream.muteVideo();
        }
        if (constraints.audio) {
          stream.unmuteAudio();
        } else {
          stream.muteAudio();
        }
        const player = <Player stream={stream} />;
        stream.playerComp = player;
        this.setState({ cameraStream: stream });
        this.props.mc?.publish(stream);
      },
      (err: Error) => {
        logger.warn('stream init error', err);
      }
    );

    callback && callback();

    stream.on('track-ended', () => {
      modalError('lock_error_track');
    });
  }

  changeShareState(s: boolean): void {
    if (s && this.props.meeting.meetingInfo.screen_shared_uid) {
      message.warn('当前有人正在分享，请等待结束后再开启');
      return;
    }
    if (this.state.screenStream) {
      this.state.screenStream?.close();
    }

    if (s) {
      const screenStream = RTC.createStream({
        screenAudio: true,
        video: true,
        screen: true,
      });

      screenStream.setVideoEncoderConfiguration(
        this.props.settings.screenStreamSettings
      );

      screenStream.init(
        () => {
          logger.debug('screen stream init succ');
          this.props.mc?.startShareScreen({});
          this.setState({ screenStream });
          this.props.mc?.publish(screenStream);

          screenStream.on('track-ended', () => {
            this.changeShareState(false);
          });

          const player = <Player stream={screenStream} />;
          screenStream.playerComp = player;
        },
        (err: Error) => {
          this.changeShareState(false);
          if (err.message === 'NotAllowedError') {
            message.error(TOASTS['screen_not_allow']);
          } else {
            message.error(TOASTS['screen_error']);
          }
        }
      );
    } else {
      this.props.mc?.endShareScreen();
      if (this.state.screenStream) {
        this.props.mc?.unpublish(this.state.screenStream as Stream);
        this.setState({ screenStream: null });
      }
    }
    this.props.setIsSharing(s);
  }

  changeUserDrawerVisibility(v: boolean): void {
    this.setState({ usersDrawerVisible: v });
  }

  recordMeeting(): void {
    const { setMeetingInfo, meeting } = this.props;
    const users = meeting.meetingUsers.map((user) => user.user_id);
    const screenUid = meeting.meetingInfo.screen_shared_uid;
    try {
      this.props.mc
        ?.recordMeeting({
          users,
          screen_uid: screenUid,
        })
        .then(() => {
          setMeetingInfo({
            ...meeting.meetingInfo,
            record: true,
          });
        });
    } catch (error) {
      message.error(`${error}`);
    }
  }

  leaving(): void {
    this.setState({ leaving: true });
  }

  get roomId(): string {
    return (
      this.props.currentUser.roomId ||
      (Utils.getQueryString('roomId') as string)
    );
  }

  end(): void {
    history.push(`/?roomId=${this.roomId}`);
  }

  closeAllStream() {
    const { cameraStream } = this.state;
    if (cameraStream) {
      //Unpublish stream when a compatible user refreshes or closes a web page
      this.props.meeting.status !== 'end' &&
        this.props.mc?.unpublish(cameraStream as Stream).then(() => {
          cameraStream?.close();
        });
    }
    if (this.state.screenStream) {
      if (this.props.currentUser.isSharing) {
        this.changeShareState(false);
      }
      this.state.screenStream?.close();
    }
  }

  //结束当前通话
  leavingMeeting(): void {
    this.props.setMeetingStatus('end');
    this.props.mc?.unpublish(this.state.cameraStream as Stream).then(() => {
      this.props.mc?.leaveMeeting({}).then(() => {
        this.setState({ leaving: false });
        this.end();
      });
    });
  }

  //结束全部通话
  endMeeting(): void {
    try {
      this.props.mc?.endMeeting({});
      this.setState({ leaving: false });
    } catch (error) {
      message.error({ content: `${error}` });
    }
  }

  changeHost(uid: string, name: string): void {
    hostChangeInfo(name, () =>
      this.props.mc
        ?.changeHost({
          user_id: uid,
        })
        .catch(() => message.error(TOASTS['give_host_error']))
    );
  }

  muteAll(): void {
    try {
      this.props.mc?.muteUser({});
    } catch (error) {
      message.error({ content: `${error}` });
    }
  }

  muteUser(user_id: string): void {
    try {
      this.props.mc?.muteUser({
        user_id,
      });
    } catch (error) {
      message.error({ content: `${error}` });
    }
  }

  askMicOn(user_id: string): void {
    try {
      this.props.mc?.askMicOn({
        user_id,
      });
      sendInfo();
    } catch (error) {
      message.error({ content: `${error}` });
    }
  }

  changeSpeakCollapse() {
    const {
      meeting: { speakCollapse },
    } = this.props;
    this.props?.setSpeakCollapse(!speakCollapse);
  }

  checkMediaState() {
    document.body.addEventListener(
      'mousemove',
      throttle(() => {
        const remoteVideoContainers = document.querySelectorAll(
          '.remote_player_container'
        );
        Array.from(remoteVideoContainers).forEach((c) => {
          const video = c.querySelector('video');
          if (video && video.muted) {
            video.muted = false;
            video.removeAttribute('muted');
          }
          if (video && video.paused) {
            video.play();
          }
        });
      }, 1000)
    );
  }

  render(): ReactNode {
    const { currentUser, meeting } = this.props;
    const {
      usersDrawerVisible,
      cameraStream,
      screenStream,
      remoteStreams,
      leaving,
      audioLevels,
      localVolume,
      refresh,
    } = this.state;
    if (refresh) {
      return <Skeleton />;
    }
    return (
      <div className={styles.container}>
        <Header
          changeSpeakCollapse={this.changeSpeakCollapse.bind(this)}
          meeting={meeting}
          username={currentUser.name || 'unknown'}
          meetingCreateAt={meeting.meetingInfo.created_at}
          now={meeting.meetingInfo.now}
          roomId={
            currentUser.roomId || meeting.meetingInfo.room_id || 'unknown'
          }
        />
        <MeetingViews
          currentUser={currentUser}
          meeting={meeting}
          cameraStream={cameraStream}
          screenStream={screenStream}
          remoteStreams={remoteStreams}
          audioLevels={audioLevels}
          localVolume={localVolume}
        />
        <ControlBar
          currentUser={currentUser}
          meeting={meeting}
          openUsersDrawer={this.changeUserDrawerVisibility.bind(this, true)}
          changeMicState={this.changeMicState.bind(this)}
          changeCameraState={this.changeCameraState.bind(this)}
          changeShareState={this.changeShareState.bind(this)}
          leaveMeeting={this.leaving.bind(this)}
          recordMeeting={this.recordMeeting.bind(this)}
        />
        <UsersDrawer
          currentUser={currentUser}
          meeting={meeting}
          visible={usersDrawerVisible}
          closeUserDrawer={this.changeUserDrawerVisibility.bind(this, false)}
          changeHost={this.changeHost.bind(this)}
          muteAll={this.muteAll.bind(this)}
          askMicOn={this.askMicOn.bind(this)}
          muteUser={this.muteUser.bind(this)}
          audioLevels={audioLevels}
          localVolume={localVolume}
        />
        <LeavingConfirm
          visible={leaving}
          isHost={currentUser.isHost}
          cancel={() => this.setState({ leaving: false })}
          leaveMeeting={this.leavingMeeting.bind(this)}
          endMeeting={this.endMeeting.bind(this)}
        />
        <StreamStats
          cameraStream={cameraStream}
          remoteStreams={remoteStreams}
        />
      </div>
    );
  }
}

export default connector(injectIntl(Meeting));
