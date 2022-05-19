import React, {
  useMemo,
  useCallback,
  useRef,
} from 'react';
import { Drawer, Tooltip, Button, message } from 'antd';
import IconBtn from '@/components/IconBtn';
import { MeetingUser } from '@/models/meeting';
import Logger from '@/utils/Logger';
import closeIcon from '/assets/images/closeIcon.png';
import shareOffIcon from '/assets/images/shareOffIcon.png';
import camOnIcon from '/assets/images/camOnIcon.png';
import camOffIcon from '/assets/images/camOffIcon.png';
import changeHostIcon from '/assets/images/changeHostIcon.png';
import muteAllIcon from '/assets/images/muteAllIcon.png';
import { connector, injectProps, ConnectedProps } from '../../configs/config';

import styles from './index.less';
import { TOASTS } from '@/constant';
import DeviceController from '@/lib/DeviceController';

import {
  sendInfo,
  hostChangeInfo,
} from '../../components/MessageTips';

interface IProps {
  visible: boolean;
  closeUserDrawer: () => void;
}

const logger = new Logger('UsersDrawer');

export type IUsersDrawerProps = ConnectedProps<typeof connector> & IProps;

const UsersDrawer: React.FC<IUsersDrawerProps> = (props) => {

  const {
    currentUser,
    mc,
    meeting,
    visible,
    closeUserDrawer,
  } = props;

  logger.debug('meeting.meetingUsers %o', meeting.meetingUsers);

  const deviceController = useRef<DeviceController>(
    new DeviceController(props)
  );

  const canMuteAll = useMemo(() => {
    return meeting.meetingUsers.filter((item) => item?.is_mic_on)?.length !== 0;
  }, [meeting.meetingUsers]);

  const commomProps = {
    width: 32,
    height: 32,
    style: { background: 'transparent', margin: 0 },
  };

  /**
   * @brief 静音所有人
   * @function muteAll
   */
  const muteAll = useCallback(() => {
    try {
      // 主持人静音除了自己以外的人
      mc?.muteUser({});
    } catch (error) {
      message.error({ content: `${error}` });
    }
  },[mc]);

  /**
   * @brief 更换主持人
   * @function changeHost
   */
  const changeHost = useCallback((uid: string, name: string): void => {
    hostChangeInfo(name, () =>
      mc
        ?.changeHost({
          user_id: uid,
        })
        .catch(() => message.error(TOASTS['give_host_error']))
    );
  },[mc]);

  /**
   * @brief 请求某用户打开麦克风
   * @function askMicOn
   */
  const askMicOn = useCallback((user_id: string): void => {
    try {
      mc?.askMicOn({
        user_id,
      });
      sendInfo();
    } catch (error) {
      message.error({ content: `${error}` });
    }
  },[mc]);

  /**
   * @brief 静音某用户
   * @function muteUser
   */
  const muteClick = useCallback((user: MeetingUser): void => {
    if (currentUser.userId === user.user_id) {
      // 如果操作的是自己
      if (deviceController.current) {
        deviceController.current.changeAudioState(!user.is_mic_on);
      }
    } else {
      if (currentUser.isHost) {
        // 请求对应用户打开 mic
        if (!user.is_mic_on) askMicOn(user.user_id);
        else {
          mc?.muteUser({
            user_id: user.user_id,
          });
        }
      } else {
        message.warn('您不是主持人, 请联系主持人进行操作');
      }
    }
  },[askMicOn, currentUser.isHost, currentUser.userId, mc]);

  /**
   * @brief 判断传入的 user_id 是否有声音
   * @function hasVolume
   */
  const hasVolume = (user_id: string) => {
    const { volumeSortList: audioLevels, localSpeaker: localUser } =
      meeting.meetingInfo;

    if (currentUser.userId && localUser) {
      // 获取 user_id 用户的音量
      const speaker = audioLevels.find(item => item.userId === user_id);
      if(user_id === currentUser.userId && currentUser.isMicOn && localUser.volume > 0.2){
        return true;
      }else{
        return speaker && speaker.volume > 0 ? true : false;
      }
    }
    return false;
  };

  /**
   * @brief 操作摄像头
   * @function onCameraChange
   */
  const onCameraChange = (user: MeetingUser):void => {
    if (currentUser.userId === user.user_id) {
      // 如果操作的是自己
      if (deviceController.current) {
        deviceController.current.changeVideoState(!user.is_camera_on);
      }
    } else {
      if (currentUser.isHost) {
        // 请求对应用户打开摄像头
        if (!user.is_camera_on) mc?.askCameraOn({user_id: user.user_id});
        else message.warn('对方摄像头已打开');
      } else {
        message.warn('您不是主持人, 请联系主持人进行操作');
      }
    }
  };

  const orderMeetingUser = useMemo(() => {
    const _users = [...meeting.meetingUsers];
    // host 在最前面, 其它顺序一样
    _users.sort(function (x, y) {
      return x.is_host ? -1 : y.is_host ? 1 : 0;
    });
    return _users;
  }, [meeting.meetingUsers]);

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
      maskClosable={true}
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
        {/* 遍历会议内所有人 */}
        {orderMeetingUser.map((user) => (
          <div className={styles.userContainer} key={user.user_id}>
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
              {/* 麦克风操作按钮 */}
              <IconBtn
                {...commomProps}
                onClick={() => {
                  muteClick(user);
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
              {/* 如果自己是主持人 && Hover 的不是自己 */}
              {currentUser.isHost && currentUser.userId !== user.user_id ? (
                <IconBtn
                  {...commomProps}
                  onClick={() => {
                    changeHost(user.user_id, user.user_name);
                  }}
                >
                  <Tooltip title="移交主持人">
                    <img src={changeHostIcon} />
                  </Tooltip>
                </IconBtn>
              ) : (
                <IconBtn
                  {...commomProps}
                  onClick={() => {
                    onCameraChange(user);
                  }}
                >
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

export default injectProps(UsersDrawer);
