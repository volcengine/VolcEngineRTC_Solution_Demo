import React, { useState, useEffect } from 'react';
import { UserModelState } from '@/models/user';
import { MeetingModelState, MeetingUser } from '@/models/meeting';
import { ViewMode, meetingActions } from '@/models/meeting';
import { Stream, AppState } from '@/app-interfaces';
import View from '@/components/View';
import Logger from '@/utils/Logger';
import { Dispatch } from '@@/plugin-dva/connect';
import { connect, bindActionCreators } from 'dva';
import { injectIntl } from 'umi';
import { ConnectedProps } from 'react-redux';
import { WrappedComponentProps } from 'react-intl';
import { IRemoteAudioLevel } from '@/app-interfaces';

import GalleryView from '../GalleryView';
import SpeakerView from '../SpeakerView';
import styles from './index.less';

interface IMeetingViewsProps {
  currentUser: UserModelState;
  meeting: MeetingModelState;
  cameraStream: Stream | null;
  screenStream: Stream | null;
  remoteStreams: { [id: string]: Stream };
  audioLevels: IRemoteAudioLevel[];
  localVolume: number;
}

const logger = new Logger('MeetingViews');
export interface ActiveMeetingUser extends MeetingUser {
  stream: Stream | null;
  me?: boolean;
  speaking?: boolean;
}

function mapStateToProps(state: AppState) {
  return {
    mc: state.meetingControl.sdk,
    // meeting: state.meeting,
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators({ ...meetingActions }, dispatch),
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export type MeetingViewsProps = ConnectedProps<typeof connector> &
  WrappedComponentProps &
  IMeetingViewsProps;


const HiddenVideo: React.FC<{
  users: MeetingUser[];
  remoteStreams: {
    [id: string]: Stream;
  };
}> = ({ users, remoteStreams }) => {
  return (
    <div className={styles['videoHidden']}>
      {users.map((item) => {
        const stream = remoteStreams[item.user_id];
        if (stream?.playerComp) {
          return stream.playerComp;
        }
        return <div key={item.user_id}></div>;
      })}
    </div>
  );
};

const MeetingViews: React.FC<MeetingViewsProps> = (props) => {
  const {
    meeting,
    currentUser,
    cameraStream,
    screenStream,
    remoteStreams,
    audioLevels,
    localVolume,
  } = props;
  const [activeUsersViews, updateActiveUsersViews] = useState<React.ReactNode[]>([]);
  const [screenView, updateScreenView] = useState<React.ReactNode>(null);

  useEffect(() => {
    if (meeting.orderMeetingUsers?.length) {
      //按声音大小来排序
      const activeUsersViews = meeting.orderMeetingUsers
        ?.slice(0, meeting.viewMode === ViewMode.GalleryView ? 9 : 8)
        .map((activeUser) => {
          const me = activeUser.user_id === currentUser.userId;
          const localProps: any = {};
          let stream: Stream | null = null;

          if (me && cameraStream) {
            stream = cameraStream;
          } else {
            stream = remoteStreams[activeUser.user_id];
          }

          logger.debug(
            'user_id: %s, stream: %o',
            activeUser.user_id,
            stream
          );

          if (me) {
            localProps.is_camera_on = currentUser.isCameraOn;
            localProps.is_mic_on = currentUser.isMicOn;
          }
          return (
            <View
              player={stream?.playerComp}
              key={activeUser.user_id}
              me={me}
              stream={stream}
              speaking={activeUser.is_mic_on}
              audioLevels={audioLevels}
              {...activeUser}
              {...localProps}
              is_sharing={false}
              sharingId={meeting.meetingInfo.screen_shared_uid}
              count={meeting.orderMeetingUsers?.length}
              localVolume={localVolume}
            />
          );
        });
      updateActiveUsersViews(activeUsersViews);
    }
  }, [meeting.viewMode, meeting.meetingUsers, cameraStream, currentUser, remoteStreams, audioLevels, meeting.meetingInfo.screen_shared_uid, localVolume, meeting.orderMeetingUsers]);

  useEffect(() => {
    if (screenStream) {
      const sharingUser = meeting.meetingUsers.find(
        (user) => user.user_id === screenStream.uid
      );
      const screenView = (
        <View
          player={screenStream?.playerComp}
          me={false}
          stream={screenStream}
          audioLevels={audioLevels}
          speaking={false}
          {...(sharingUser as MeetingUser)}
          is_sharing={true}
          localVolume={localVolume}
        />
      );
      updateScreenView(screenView);
    }
  }, [
    meeting.meetingUsers,
    screenStream,
    currentUser,
    audioLevels,
    localVolume,
  ]);

  const renderViews = () => {
    if (meeting.viewMode === ViewMode.GalleryView) {
      return (
        <>
          <GalleryView views={activeUsersViews} />
          <HiddenVideo
            users={meeting.orderMeetingUsers.slice(
              9,
              meeting.orderMeetingUsers.length
            )}
            remoteStreams={remoteStreams}
          />
        </>
      );
    } else {
      return (
        <>
          <SpeakerView
            screenView={screenView}
            views={activeUsersViews}
            meeting={meeting}
          />
          <HiddenVideo
            users={meeting.orderMeetingUsers.slice(
              8,
              meeting.orderMeetingUsers.length
            )}
            remoteStreams={remoteStreams}
          />
        </>
      );
    }
  };

  return <div className={styles.container}>{renderViews()}</div>;
};

export default connector(injectIntl(MeetingViews));
