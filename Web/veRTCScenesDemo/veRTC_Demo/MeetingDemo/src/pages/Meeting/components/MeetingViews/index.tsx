import React, { useState, useEffect } from 'react';
import { UserModelState } from '@/models/user';
import { MeetingUser } from '@/models/meeting';
import { ViewMode } from '@/models/meeting';
import { Stream } from '@/app-interfaces';
import View from '@/components/View';
import Logger from '@/utils/Logger';
import { WrappedComponentProps } from 'react-intl';
import { injectProps, ConnectedProps, connector } from '../../configs/config';
import { LocalPlayer, ShareView } from '../../components/MediaPlayer';

import GalleryView from '../GalleryView';
import SpeakerView from '../SpeakerView';
import styles from './index.less';

interface IMeetingViewsProps {
  currentUser: UserModelState;
  cameraStream: Stream | null;
  screenStream: boolean;
  remoteStreams: { [id: string]: Stream };
}

const logger = new Logger('MeetingViews');
export interface ActiveMeetingUser extends MeetingUser {
  me?: boolean;
  speaking?: boolean;
}

const HiddenVideo: React.FC<{
  views: React.ReactNode[];
}> = ({ views }) => {
  return (
    <div className={styles['videoHidden']}>
      {views.map((view) => {
        return <>{view}</>;
      })}
    </div>
  );
};

export type MeetingViewsProps = ConnectedProps<typeof connector> &
  WrappedComponentProps &
  IMeetingViewsProps;


const MeetingViews: React.FC<MeetingViewsProps> = (props) => {
  const {
    meeting,
    currentUser,
    cameraStream,
    screenStream,
    remoteStreams,
  } = props;
  const [activeUsersViews, updateActiveUsersViews] = useState<React.ReactNode[]>([]);
  const [screenView, updateScreenView] = useState<React.ReactNode>(null);
  const [localView, updateLocalView] = useState<React.ReactNode>(null);

  useEffect(() => {
    const volumeSortList = meeting.meetingInfo.volumeSortList;
    if (volumeSortList?.length) {
      //按声音大小来排序
      const activeUsersViews = volumeSortList
        .map(({ userId, volume }) => {
          const localProps: any = {};
          const activeUser = meeting.meetingUsers.find(
            (i) => i.user_id === `${userId}`
          );
          return activeUser ? (
            <View
              player={remoteStreams[userId]?.playerComp}
              key={userId}
              me={false}
              speaking={activeUser?.is_mic_on || false}
              volume={volume ?? 0}
              {...activeUser}
              {...localProps}
              is_sharing={false}
              sharingId={meeting.meetingInfo.screen_shared_uid}
              count={volumeSortList?.length}
              sharingView={false}
            />
          ) : null;
        });
      updateActiveUsersViews(activeUsersViews.filter((i) => i));
    }
  }, [
    meeting.viewMode,
    meeting.meetingUsers,
    cameraStream,
    currentUser,
    remoteStreams,
    meeting.meetingInfo.screen_shared_uid,
    meeting.meetingInfo
  ]);

  //分享流View
  useEffect(() => {
    if (
      meeting.viewMode === ViewMode.SpeakerView &&
      meeting.meetingInfo.screen_shared_uid &&
      meeting.isSharing
    ) {
      const sharingUser = meeting.meetingUsers.find(
        (user) => user.user_id === meeting.meetingInfo.screen_shared_uid
      );
      const me = sharingUser?.user_id === currentUser.userId;
      const screenView = (
        <View
          player={
            <ShareView
              rtc={props.rtc}
              sharingUser={sharingUser}
              me={me}
              localCaptureSuccess={meeting.localCaptureSuccess}
            />
          }
          me={me}
          speaking={sharingUser?.is_mic_on}
          {...(sharingUser as MeetingUser)}
          is_sharing={true}
          sharingView={true}
          volume={meeting.meetingInfo.localSpeaker?.volume ?? 0}
        />
      );
      updateScreenView(screenView);
    }
  }, [
    meeting.isSharing,
    meeting.localCaptureSuccess,
    meeting.meetingUsers,
    screenStream,
    currentUser,
    props.rtc,
    meeting.meetingInfo.screen_shared_uid,
    meeting.viewMode,
    meeting.meetingInfo.localSpeaker.volume,
  ]);

  //本地流View
  useEffect(() => {
    const _user = meeting.meetingUsers.find(
      (user) => user.user_id === currentUser.userId
    );
    const localView = (
      <View
        player={
          <LocalPlayer
            localCaptureSuccess={meeting.localCaptureSuccess}
            rtc={props.rtc}
            renderDom="local-player"
          />
        }
        key={currentUser.userId}
        me={true}
        speaking={currentUser?.isMicOn || false}
        sharingId={meeting.meetingInfo.screen_shared_uid}
        volume={meeting.meetingInfo.localSpeaker?.volume ?? 0}
        sharingView={false}
        {...(_user as MeetingUser)}
        {...{
          is_camera_on: currentUser.isCameraOn,
          is_mic_on: currentUser.isMicOn,
        }}
      />
    );
    updateLocalView(localView);
  }, [
    meeting.meetingInfo.localSpeaker.volume,
    meeting.meetingUsers,
    currentUser,
    meeting.localCaptureSuccess,
    meeting.meetingInfo.screen_shared_uid,
    props.rtc,
  ]);


  //宫格模式
  if (meeting.viewMode === ViewMode.GalleryView) {
    const galleryUser = activeUsersViews.slice(0, 8);
    return (
      <div className={styles.container}>
        <GalleryView views={[localView, ...galleryUser]} />
        <HiddenVideo
          views={activeUsersViews.slice(8, activeUsersViews.length)}
        />
      </div>
    );
  }

  //分享模式
  return (
    <div className={styles.container}>
      <SpeakerView
        screenView={screenView}
        views={[localView, ...activeUsersViews.slice(0, 7)]}
        meeting={meeting}
      />
      <HiddenVideo views={activeUsersViews.slice(7, activeUsersViews.length)} />
    </div>
  );
};

export default injectProps(MeetingViews);
