import React, { useCallback, useRef, useState } from 'react';
import { Tooltip, message } from 'antd';
import IconBtn from '@/components/IconBtn';
import SettingsModal from '@/components/SettingsModal';
import { modalWarning } from '@/pages/Meeting/components/MessageTips';
import styles from './index.less';
import micOnIcon from '/assets/images/micOnIcon.png';
import micOffIcon from '/assets/images/micOffIcon.png';
import camOnIcon from '/assets/images/camOnIconMeeting.png';
import camOffIcon from '/assets/images/camOffIcon.png';
import shareOnIcon from '/assets/images/shareOnIcon.png';
import shareOffIcon from '/assets/images/shareOffIcon.png';
import recordOnIcon from '/assets/images/recordOnIcon.png';
import recordOffIcon from '/assets/images/recordOffIcon.png';
import usersIcon from '/assets/images/usersIcon.png';
import settingIcon from '/assets/images/settingIcon.png';
import endIcon from '/assets/images/endIcon.png';
import {
  connector,
  injectProps,
  ConnectedProps,
} from '@/pages/Meeting/configs/config';
import DeviceController from '@/lib/DeviceController';

interface IControlBarProps {
  openUsersDrawer: () => void;
  leaveMeeting: () => void;
}

export type ControlBarProps = ConnectedProps<typeof connector> &
  IControlBarProps;

const ControlBar: React.FC<ControlBarProps> = (props) => {
  const { currentUser, meeting, openUsersDrawer, leaveMeeting } = props;

  const deviceController = useRef<DeviceController>(
    new DeviceController(props)
  );

  const commonProps = {
    width: 36,
    height: 36,
    style: { background: 'transparent', margin: 0 },
  };

  /**
   * @param visible 设置窗口是否可见
   */
  const [visible, setVisible] = useState(false);

  /**
   * @brief 麦克风切换状态
   * @function changeMicState
   */
  const changeMicState = useCallback((micState: boolean): void => {
    if (deviceController.current) {
      deviceController.current.changeAudioState(micState);
    }
  }, []);

  /**
   * @brief 摄像头切换状态
   * @function changeCameraState
   */
  const changeCameraState = useCallback((cameraState: boolean): void => {
    if (deviceController.current) {
      deviceController.current.changeVideoState(cameraState);
    }
  }, []);

  const changeShareState = (isShare: boolean) => {
    const { meeting, settings } = props;
    const param = {
      meeting,
      settings,
    };
    deviceController?.current.changeShareState(param, isShare);
  };

  /**
   * @brief 会议录像
   * @function recordMeeting
   */
  const recordMeeting = (): void => {
    const { setMeetingInfo, meeting } = props;
    const _users = meeting.meetingUsers.map((user) => user.user_id);
    const screenUid = meeting.meetingInfo.screen_shared_uid;
    try {
      props.mc
        ?.recordMeeting({
          users: _users,
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
  };

  return (
    <div className={styles.container}>
      <div className={styles.funcBtn}>
        <IconBtn
          {...commonProps}
          onClick={() => {
            const { audio, audioMessage } = currentUser.deviceAccess;
            if (!audio) {
              audioMessage && modalWarning(audioMessage);
              return;
            }
            changeMicState(!currentUser.isMicOn);
          }}
        >
          <Tooltip title="麦克风">
            <img src={currentUser.isMicOn ? micOnIcon : micOffIcon} />
          </Tooltip>
        </IconBtn>
      </div>

      <div className={styles.funcBtn}>
        <IconBtn
          {...commonProps}
          onClick={() => {
            const { video, videoMessage } = currentUser.deviceAccess;
            if (!video) {
              videoMessage && modalWarning(videoMessage);
              return;
            }
            changeCameraState(!currentUser.isCameraOn);
          }}
        >
          <Tooltip title="摄像头">
            <img src={currentUser.isCameraOn ? camOnIcon : camOffIcon} />
          </Tooltip>
        </IconBtn>
      </div>

      <div className={styles.funcBtn}>
        <IconBtn
          {...commonProps}
          onClick={() => changeShareState(!currentUser.isSharing)}
        >
          <Tooltip title="屏幕共享">
            <img
              src={
                !meeting.meetingInfo.screen_shared_uid
                  ? shareOnIcon
                  : shareOffIcon
              }
            />
          </Tooltip>
        </IconBtn>
      </div>

      <div className={styles.funcBtn}>
        <IconBtn
          {...commonProps}
          onClick={() => {
            if (!meeting.meetingInfo.record) {
              recordMeeting();
            }
          }}
        >
          <Tooltip
            title={
              !meeting.meetingInfo.record ? '开启录制' : '暂不支持停止录制'
            }
          >
            <img
              src={meeting.meetingInfo.record ? recordOffIcon : recordOnIcon}
            />
          </Tooltip>
        </IconBtn>
      </div>

      <div
        className={styles.funcBtn}
        onClick={() => {
          if (currentUser.isSharing) {
            message.info('停止屏幕共享后可查看参会者列表');
            return;
          }
          openUsersDrawer();
        }}
      >
        <IconBtn {...commonProps}>
          <Tooltip title="参会人列表">
            <img src={usersIcon} />
          </Tooltip>
        </IconBtn>
      </div>

      <div className={styles.funcBtn}>
        <IconBtn
          {...commonProps}
          onClick={() => {
            if (currentUser.isSharing) {
              message.info('停止屏幕共享后可进入会议设置页');
              return;
            }
            setVisible(true);
          }}
        >
          <Tooltip title="设置">
            <img src={settingIcon} />
          </Tooltip>
        </IconBtn>
      </div>

      <span className={styles.split} />

      <div className={styles.funcBtn}>
        <IconBtn {...commonProps} onClick={leaveMeeting}>
          <Tooltip title="结束通话">
            <img src={endIcon} />
          </Tooltip>
        </IconBtn>
      </div>
      {/* 设置 Modal */}
      <SettingsModal visible={visible} close={() => setVisible(false)} />
    </div>
  );
};

export default injectProps(ControlBar);
