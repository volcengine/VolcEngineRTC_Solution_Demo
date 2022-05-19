import React, { Component } from 'react';
import { message } from 'antd';
import { history } from 'umi';
import Logger from '@/utils/Logger';
import Utils from '@/utils/utils';
import { ViewMode } from '@/models/meeting';
import { UsersDrawer, LeavingConfirm, ControlBar, Header } from './components';
import { MeetingProps } from '@/app-interfaces';
import MeetingEvent from '@/pages/Meeting/components/MeetEvent';
import VideoAudioSubscribe from '@/lib/VideoAudioSubscribe';
import { injectProps } from '@/pages/Meeting/configs/config';
import { modalError } from './components/MessageTips';

import styles from './index.less';

const logger = new Logger('meeting');

/**
 * @param drawerVisible 用户列表抽屉是否可见
 * @param exitVisible 退出会议弹窗是否可见
 */
const initState = {
  drawerVisible: false,
  exitVisible: false,
};

type MeetingState = {
  drawerVisible: boolean;
  exitVisible: boolean;
};

class Meeting extends Component<MeetingProps, MeetingState> {
  constructor(props: MeetingProps) {
    super(props);

    window.__meeting = this;
    this.state = initState;
    this.leavingMeeting = this.leavingMeeting.bind(this);
    this.endMeeting = this.endMeeting.bind(this);
    this.unMount = this.unMount.bind(this);
  }

  // deviceLib = new DeviceController(this.props);
  subscribeLib = new VideoAudioSubscribe(this.props);

  get roomId(): string {
    return (
      this.props.currentUser.roomId ||
      (Utils.getQueryString('roomId') as string)
    );
  }

  componentDidMount = (): void => {
    Utils.checkMediaState();
    window.addEventListener('beforeunload', this.unMount);
  };

  componentWillUnmount() {
    this.unMount();
    window.removeEventListener('beforeunload', this.unMount);
  }

  componentDidUpdate(prevProps: MeetingProps, preState: MeetingState): void {
    if (
      this.props.mc &&
      prevProps.currentUser.isHost !== this.props.currentUser.isHost
    ) {
      this.props.mc.isHost = this.props.currentUser.isHost;
    }
  }

  unMount() {
    const props = this.props;
    props.rtc.leave();
    props.rtc.removeEventListener();
    props.mc?.removeEvent();
    this.subscribeLib.cleanStreamRecord();
    // this.deviceLib?.stopShare(false);
    props.setMeetingUsers([]);
    props.setLocalCaptureSuccess(false);
    props.setViewMode(ViewMode.GalleryView);
    props.setSpeakCollapse(false);
    this.setState({
      ...initState,
    });
  }

  end(): void {
    history.push(`/?roomId=${this.roomId}`);
  }

  // //结束当前通话
  leavingMeeting(): void {
    const { props } = this;
    props.setMeetingStatus('end');
    props.rtc.leave();
    props.mc?.leaveMeeting({});
    if (props.meeting.meetingInfo.screen_shared_uid) {
      props.mc?.endShareScreen();
    }
    props.rtc.removeEventListener();
    // this.setState({ leaving: false });
    this.end();
  }

  // //结束全部通话
  endMeeting(): void {
    try {
      this.props.mc?.endMeeting({});
      // this.setState({ leaving: false });
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

  render() {
    const { props, state } = this;

    return (
      <div className={styles.container}>
        <Header
          changeSpeakCollapse={this.changeSpeakCollapse.bind(this)}
          meeting={props.meeting}
          username={props.currentUser.name || 'unknown'}
          meetingCreateAt={props.meeting.meetingInfo.created_at}
          now={props.meeting.meetingInfo.now}
          roomId={
            props.currentUser.roomId ||
            props.meeting.meetingInfo.room_id ||
            'unknown'
          }
        />

        <MeetingEvent leavingMeeting={this.leavingMeeting} end={this.end} />

        <ControlBar
          openUsersDrawer={() => this.setState({ drawerVisible: true })}
          leaveMeeting={() => this.setState({ exitVisible: true })}
        />

        <UsersDrawer
          visible={state.drawerVisible}
          closeUserDrawer={() => this.setState({ drawerVisible: false })}
        />

        <LeavingConfirm
          visible={state.exitVisible}
          isHost={props.currentUser.isHost}
          cancel={() => {
            this.setState({ exitVisible: false });
          }}
          leaveMeeting={this.leavingMeeting}
          endMeeting={this.endMeeting}
        />
      </div>
    );
  }
}

export default injectProps(Meeting);
