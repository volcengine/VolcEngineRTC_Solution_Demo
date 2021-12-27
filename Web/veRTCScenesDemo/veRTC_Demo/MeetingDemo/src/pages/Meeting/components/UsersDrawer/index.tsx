import React, { useState, useMemo } from 'react';
import { Drawer, Tooltip, Button } from 'antd';
import IconBtn from '@/components/IconBtn';
import { UserModelState } from '@/models/user';
import { MeetingModelState } from '@/models/meeting';
import Logger from '@/utils/Logger';
import { IRemoteAudioLevel } from '@/app-interfaces';
import closeIcon from '/assets/images/closeIcon.png';
import shareOffIcon from '/assets/images/shareOffIcon.png';
import camOnIcon from '/assets/images/camOnIcon.png';
import camOffIcon from '/assets/images/camOffIcon.png';
import changeHostIcon from '/assets/images/changeHostIcon.png';
import muteAllIcon from '/assets/images/muteAllIcon.png';

import styles from './index.less';

interface IUsersDrawerProps {
  currentUser: UserModelState;
  meeting: MeetingModelState;
  visible: boolean;
  closeUserDrawer: () => void;
  changeHost: (uid: string, name: string) => void;
  muteAll: () => void;
  askMicOn: (user_id: string) => void;
  muteUser: (user_id: string) => void;
  audioLevels: IRemoteAudioLevel[];
  localVolume: number;
}

const logger = new Logger('UsersDrawer');

const UsersDrawer: React.FC<IUsersDrawerProps> = ({
  currentUser,
  meeting,
  visible,
  closeUserDrawer,
  changeHost,
  muteAll,
  askMicOn,
  muteUser,
  audioLevels,
  localVolume,
}) => {
  logger.debug('meeting.meetingUsers %o', meeting.meetingUsers);

  const [hoverUserId, updateHoverUserId] = useState<string>('');

  const canMuteAll = useMemo(() => {
    return meeting.meetingUsers.filter((item) => item?.is_mic_on)?.length !== 0;
  }, [meeting.meetingUsers]);

  const commomProps = {
    width: 32,
    height: 32,
    style: { background: 'transparent', margin: 0 },
  };

  const hasVolume = (user_id: string) => {
    if(user_id === currentUser.userId && currentUser.isMicOn && localVolume > 0.2){
      return true;
    }else{
      const remoteAudio = audioLevels.find((item) => item?.userId === user_id);
      if(remoteAudio && remoteAudio?.RecvLevel > 0){
        return true;
      }else{
        return false;
      }
    }
  };

  const orderMeetingUser = useMemo(() => {
    const _users = [...meeting.meetingUsers];
    _users.sort(function (x, y) {
      return x.is_host ? -1 : y.is_host ? 1 : 0;
    });
    return _users;
  },[meeting.meetingUsers]);

  return (
    <Drawer
      className={styles.container}
      visible={visible}
      width={280}
      bodyStyle={{
        position: 'relative',
        backgroundColor: '#101319',
        color: '#fff',
        padding: 0,
        fontSize: 12,
      }}
      headerStyle={{ backgroundColor: '#101319', border: 'none' }}
      title={
        <h3 style={{ color: '#fff', fontSize: 16, marginBottom: 0 }}>
          参会人（{meeting.meetingUsers.length}）
        </h3>
      }
      closable={false}
      extra={
        <img
          style={{ width: 13, cursor: 'pointer' }}
          src={closeIcon}
          onClick={closeUserDrawer}
        />
      }
    >
      <div className={styles.useContainerWrapper}>
        {orderMeetingUser.map((user) => (
          <div
            className={styles.userContainer}
            key={user.user_id}
            onMouseEnter={() => {
              if (currentUser.isHost && currentUser.userId !== user.user_id) {
                updateHoverUserId(user.user_id);
              }
            }}
            onMouseLeave={() => updateHoverUserId('')}
          >
            <div className={styles.left}>
              <span
                className={styles.avatar}
                style={{
                  border:
                    hasVolume(user.user_id) && user.is_mic_on
                      ? '1px solid #3370FF'
                      : 'none',
                }}
              >
                {user.user_name[0]}
              </span>
              <div className={styles.usernameContainer}>
                <span className={styles.username}>
                  <span>{user.user_name}</span>
                  {meeting.meetingInfo.screen_shared_uid === user.user_id && (
                    <IconBtn
                      width={24}
                      height={24}
                      style={{ background: 'transparent', margin: 0 }}
                    >
                      <img src={shareOffIcon} />
                    </IconBtn>
                  )}
                </span>
                <span>
                  {user.is_host && <span className={styles.tag}>主持人</span>}
                  {user.user_id === currentUser.userId && (
                    <span className={styles.tag}>我</span>
                  )}
                </span>
              </div>
            </div>
            <div className={styles.right}>
              <IconBtn
                {...commomProps}
                onClick={() => {
                  if (user.is_host) {
                    return;
                  }
                  user.is_mic_on
                    ? muteUser(user.user_id)
                    : askMicOn(user.user_id);
                }}
              >
                <Tooltip
                  title={
                    user.user_id !== currentUser.userId
                      ? user.is_mic_on
                        ? '静音'
                        : '打开麦克风'
                      : ''
                  }
                >
                  <span
                    className={
                      user.is_mic_on
                        ? hasVolume(user.user_id)
                          ? styles['micHasVolume']
                          : styles['micOnIcon']
                        : styles['micOffIcon']
                    }
                  ></span>
                </Tooltip>
              </IconBtn>
              {hoverUserId === user.user_id ? (
                <IconBtn
                  {...commomProps}
                  onClick={() => {
                    changeHost(user.user_id, user.user_name);
                    updateHoverUserId('');
                  }}
                >
                  <Tooltip title="移交主持人">
                    <img src={changeHostIcon} />
                  </Tooltip>
                </IconBtn>
              ) : (
                <IconBtn {...commomProps}>
                  <img src={user.is_camera_on ? camOnIcon : camOffIcon} />
                </IconBtn>
              )}
            </div>
          </div>
        ))}
      </div>
      {currentUser.isHost && (
        <div className={styles.useMuteAll}>
          <Button
            className={styles.useMuteAllButton}
            disabled={!canMuteAll}
            onClick={muteAll}
          >
            <span>
              <img src={muteAllIcon} />
            </span>
            全体静音
          </Button>
        </div>
      )}
    </Drawer>
  );
};

export default UsersDrawer;
