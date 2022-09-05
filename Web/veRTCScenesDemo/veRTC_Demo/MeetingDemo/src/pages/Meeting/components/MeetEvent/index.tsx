import { message, Skeleton } from 'antd';
import Logger from '@/utils/Logger';
import { TOASTS } from '@/config';

import type {
  UserStatusChangePayload,
  HostChangePayload,
  JoinMeetingResponse,
} from '@/lib/socket-interfaces';
import { ViewMode, MeetingUser } from '@/models/meeting';
import {
  showOpenMicConfirm,
  sendMutedInfo,
} from '@/pages/Meeting/components/MessageTips';
import React from 'react';
import DeviceController from '@/lib/DeviceController';
import { injectProps } from '@/pages/Meeting/configs/config';
import Utils from '@/utils/utils';
import { history } from 'umi';
import VideoAudioSubscribe from '@/lib/VideoAudioSubscribe';
import { IMeetingState, IVolume, MeetingProps } from '@/app-interfaces';
import {
  RTCStream,
  RemoteStreamStats,
  LocalStreamStats,
} from '@volcengine/rtc';
import { MediaPlayer } from '../MediaPlayer';
import MeetingViews from '../MeetingViews';
import StreamStats from '../StreamStats';

const logger = new Logger('meeting-event');

interface IEvent {
  end: () => void;
  leavingMeeting: () => void;
}

type IProps = MeetingProps & IEvent;

/**
 * @param drawerVisible 用户列表抽屉是否可见
 * @param exitVisible 退出会议弹窗是否可见
 * @param volumeSortList 音量根据大到小排序的对象数组
 */
const initState = {
  usersDrawerVisible: false,
  cameraStream: null,
  screenStream: false,
  remoteStreams: {},
  leaving: false,
  audioLevels: [],
  localVolume: 0,
  refresh: false,
  volumeSortList: [],
  localSpeaker: {
    userId: '',
    volume: 0,
  },
  streamStatses: {
    local: {},
    localScreen: undefined,
    remoteStreams: {},
  },
  users: [],
};

class MeetingEvent extends React.Component<IProps, IMeetingState> {
  constructor(props: IProps) {
    super(props);

    window.__meeting = this;

    this.state = initState;

    this.handleStreamAdd = this.handleStreamAdd.bind(this);
    this.handleStreamRemove = this.handleStreamRemove.bind(this);
    this.handleEventError = this.handleEventError.bind(this);
    this.handleAudioVolumeIndication =
      this.handleAudioVolumeIndication.bind(this);
    this.handleLocalStreamState = this.handleLocalStreamState.bind(this);
    this.handleRemoteStreamState = this.handleRemoteStreamState.bind(this);
    this.hanldleMeetingController = this.hanldleMeetingController.bind(this);
    this.handleTrackEnd = this.handleTrackEnd.bind(this);
    this.hanldleMeetingController();
  }

  deviceLib = new DeviceController(this.props);
  subscribeLib = new VideoAudioSubscribe(this.props);

  get roomId(): string {
    return (
      this.props.currentUser.roomId ||
      (Utils.getQueryString('roomId') as string)
    );
  }

  componentDidMount = (): void => {
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

    this.initRTC();

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

    history.listen((location, action) => {
      if (action === 'POP') {
        if (location.pathname === '/meeting') {
          this.setState({
            refresh: true,
          });
        }
      }
    });
  };

  initRTC() {
    //RTC初始化、事件绑定
    if (!this.props.rtc?.engine) {
      this.props.rtc?.createEngine();
    }
    this.props.rtc?.bindEngineEvents({
      handleStreamAdd: this.handleStreamAdd,
      handleStreamRemove: this.handleStreamRemove,
      handleEventError: this.handleEventError,
      handleAudioVolumeIndication: this.handleAudioVolumeIndication,
      handleLocalStreamState: this.handleLocalStreamState,
      handleRemoteStreamState: this.handleRemoteStreamState,
      handleTrackEnd: this.handleTrackEnd,
    });
  }

  changeUserStatus = (
    state: 'is_mic_on' | 'is_camera_on' | 'is_sharing',
    payload: UserStatusChangePayload
  ) => {
    logger.debug(payload);
    const users = this.props.meeting.meetingUsers.map((user) => {
      const ret = { ...user };
      if (ret.user_id === payload.user_id) {
        ret[state] = payload.status;
      }
      return ret;
    });
    this.props.setMeetingUsers(users);
  };

  hanldleMeetingController = () => {
    this.props.mc?.on('onUserJoinMeeting', (payload: MeetingUser) => {
      const users = this.props.meeting.meetingUsers;
      if (users.some((user) => user.user_id === payload.user_id)) {
        logger.error('onUserJoinMeeting 该用户 %s 已存在', payload.user_id);
      } else {
        this.props.setMeetingUsers([...users, payload]);
      }
    });

    // broadcasting / 广播消息
    this.props.mc?.on(
      'onUserMicStatusChange',
      (payload: UserStatusChangePayload) => {
        this.changeUserStatus('is_mic_on', payload);
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
      this.props.setMeetingUsers(users);
      const isHost = payload.host_id === this.props.currentUser.userId;
      if (isHost) {
        message.info('你已成为主持人');
      }
      this.props.setIsHost(isHost);
      this.props.setMeetingInfo({
        ...this.props.meeting.meetingInfo,
        host_id: payload.host_id,
      });
    });

    this.props.mc?.on('onUserLeaveMeeting', (payload: { user_id: string }) => {
      const users = this.props.meeting.meetingUsers.filter(
        (user) => user.user_id !== payload.user_id
      );
      this.props.setMeetingUsers([...users]);
    });

    this.props.mc?.on(
      'onShareScreenStatusChanged',
      (payload: UserStatusChangePayload) => {
        if (payload.status) {
          this.props.setViewMode(ViewMode.SpeakerView);
          this.props.setMeetingInfo({
            ...this.props.meeting.meetingInfo,
            screen_shared_uid: payload.user_id,
          });
        } else {
          this.props.setViewMode(ViewMode.GalleryView);
          this.props.setMeetingInfo({
            ...this.props.meeting.meetingInfo,
            screen_shared_uid: '',
          });
        }
        this.changeUserStatus('is_sharing', payload);
      }
    );

    this.props.mc?.on('onRecord', () => {
      this.props.setMeetingInfo({
        ...this.props.meeting.meetingInfo,
        record: true,
      });
    });

    this.props.mc?.on('onMeetingEnd', () => {
      logger.debug('onMeetingEnd');
      this.props.setMeetingStatus('end');
      this.props.rtc?.unpublish().then(() => {
        this.props.leavingMeeting();
      });
    });

    this.props.mc?.on('onMuteAll', () => {
      //user who is not host will receive muted message
      if (!this.props.currentUser.isHost) {
        this.deviceLib.changeAudioState(false);
        sendMutedInfo();
      }
      //重新设置会议人员的状态
      if (this.props.currentUser?.userId) {
        this.props.mc
          ?.getMeetingUserInfo({
            user_id: this.props.currentUser?.userId,
          })
          .then((res: any) => {
            this.props.setMeetingUsers(res);
          });
      }
    });

    this.props.mc?.on('onMuteUser', () => {
      sendMutedInfo();
      this.deviceLib.changeAudioState(false);
    });

    this.props.mc?.on('onAskingMicOn', () => {
      showOpenMicConfirm(() => this.deviceLib.changeAudioState(true));
    });

    this.props.mc?.on('onUserKickedOff', () => {
      message.warning(TOASTS['tick']);
      this.props.end();
    });

    this.props.mc?.on('onUserReconnect', (code: number) => {
      if (code === 200) {
        if (this.props.currentUser?.userId) {
          this.props.mc
            ?.getMeetingInfo()
            .then(() => {
              this.props.mc
                ?.getMeetingUserInfo({})
                .then((res: any) => {
                  this.props.setMeetingUsers(res);
                })
                .catch((err: any) => {
                  if (err === 'record not found') {
                    message.warning('会议已结束', undefined, () => {
                      this.props.end();
                    });
                  }
                });
            })
            .catch((err: any) => {
              if (err === 'record not found') {
                message.warning('会议已结束', undefined, () => {
                  this.props.end();
                });
              }
            });
        }
      } else {
        if (code === 404) {
          message.warning('你已被踢出房间', undefined, () => {
            this.props.end();
          });
          return;
        }
        if (code === 422) {
          message.warning('会议已结束', undefined, () => {
            this.props.end();
          });
          return;
        }
      }
    });
  };

  handleMeetingStatus() {
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.props.setMeetingStatus('hidden');
      }
    });
  }

  handleStreamAdd({ stream }: { stream: RTCStream }) {
    const { isScreen, userId } = stream;

    if (isScreen) {
      this.props.setMeetingIsSharing(true);
      return;
    }
    (stream as RTCStream & { playerComp: any }).playerComp = (
      <MediaPlayer
        userId={userId}
        stream={stream}
        setRemoteVideoPlayer={this.props.rtc?.setRemoteVideoPlayer}
      />
    );
    const _s = {
      ...this.state.remoteStreams,
      [userId]: stream,
    };

    this.setState({
      remoteStreams: {
        ..._s,
      },
    });
    this.subscribeLib.streams = _s;
    this.subscribeLib.addSubscribed(userId);
  }

  handleStreamRemove({ stream }: { stream: RTCStream }) {
    const { isScreen, userId } = stream;
    if (isScreen) {
      this.props.setMeetingIsSharing(false);
      return;
    }
    const remoteStreams = this.state.remoteStreams;
    if (remoteStreams[userId]) {
      delete remoteStreams[userId];
    }
    this.setState({
      remoteStreams: { ...remoteStreams },
    });
    this.subscribeLib.streams = remoteStreams;
  }

  handleEventError(e: any, VERTC: any) {
    if (e.errorCode === VERTC.ErrorCode.DUPLICATE_LOGIN) {
      message.error('你的账号被其他人顶下线了');
      //leaveRoom();
    }
  }

  handleAudioVolumeIndication(event: { speakers: IVolume[] }) {
    const otherSpeakers: IVolume[] = [];
    const meetInfo = this.props.meeting.meetingInfo;
    const host_id = meetInfo.host_id;
    for (const speaker of event.speakers) {
      if (speaker.userId !== this.props.currentUser.userId)
        otherSpeakers.push(speaker);
      else {
        this.props.setMeetingInfo({ ...meetInfo, localSpeaker: speaker });
      }
    }
    const sortList = otherSpeakers.sort((a: IVolume, b: IVolume) => {
      if (a.userId === host_id) {
        // a before b
        return -1;
      }
      if (b.userId === host_id) {
        //b before a
        return 1;
      }
      return b.volume - a.volume;
    });
    this.subscribeLib.changSubscribe(sortList);
    this.props.setMeetingInfo({ ...meetInfo, volumeSortList: sortList });
  }

  handleLocalStreamState(stats: LocalStreamStats) {
    const { streamStatses } = this.state;
    if (stats.isScreen) {
      streamStatses.localScreen = stats;
    } else {
      streamStatses.local = stats;
    }
    this.setState({
      streamStatses,
    });
  }

  handleRemoteStreamState(stats: RemoteStreamStats) {
    if (!stats || !stats.userId) return;
    const { streamStatses } = this.state;
    const key = `${stats.userId}-${stats.isScreen ? 'screen' : 'video'}`;
    streamStatses.remoteStreams[key] = stats;
    this.setState({
      streamStatses,
    });
  }

  handleTrackEnd(event: { kind: string; isScreen: boolean }) {
    const { kind, isScreen } = event;
    if (isScreen) {
      this.deviceLib.stopShare(false);
    } else {
      if (kind === 'video') {
        if (this.props.currentUser.isCameraOn) {
          this.deviceLib.changeVideoState(false);
        }
      }
    }
  }

  joinMeeting() {
    const {
      currentUser,
      setMeetingInfo,
      setMeetingUsers,
      setIsHost,
      setViewMode,
      settings,
      meeting: { meetingInfo },
    } = this.props;

    const userId = Utils.getLoginUserId();
    this.props.setUserId(userId);

    //TODO增加兜底
    if (!this.props.currentUser.appId) {
      return;
    }
    const rtc = this.props.rtc;

    const param = {
      currentUser,
      settings,
    };

    this.props.mc
      ?.joinMeeting({
        app_id: rtc.config.appId,
        user_id: userId,
        user_name: Utils.getLoginUserName(),
        room_id: this.roomId,
        mic: currentUser.isMicOn,
        camera: currentUser.isCameraOn,
      })
      .then((res?: JoinMeetingResponse) => {
        if (!res) {
          return;
        }
        const { info, users } = res;
        setMeetingInfo({ ...meetingInfo, ...info });
        const _users = this.props.meeting.meetingUsers;
        setMeetingUsers([..._users, ...users]);
        if (info.screen_shared_uid) {
          setViewMode(ViewMode.SpeakerView);
        }
        setIsHost(info.host_id === Utils.getLoginUserId());

        rtc.join(res.token, this.roomId, userId).then(() => {
          if (this.props.meeting.localAudioVideoCaptureSuccess) {
            rtc.engine.publish();
            this.props.setLocalCaptureSuccess(true);
          } else {
            this.deviceLib.openCamera(
              param,
              () => {
                this.props.setLocalCaptureSuccess(true);
              },
              true
            );
          }
          // 监听音量变化
          this.props.rtc?.engine.setAudioVolumeIndicationInterval(1000 * 1);
        });
      })
      .catch((e) => {
        message.error(`Join meeting failed: ${e || 'unknown'}`, () => {
          if (e === 'login token expired') {
            Utils.removeLoginInfo();
            history.push('/login');
          } else {
            this.props.end();
          }
        });
      });
  }

  render() {
    const { state, props } = this;
    if (state.refresh) {
      return <Skeleton />;
    }
    return (
      <>
        <MeetingViews
          currentUser={props.currentUser}
          meeting={props.meeting}
          cameraStream={state.cameraStream}
          screenStream={state.screenStream}
          remoteStreams={state.remoteStreams}
        />
        <StreamStats streamStatses={state.streamStatses} />
      </>
    );
  }
}

export default injectProps(MeetingEvent);
