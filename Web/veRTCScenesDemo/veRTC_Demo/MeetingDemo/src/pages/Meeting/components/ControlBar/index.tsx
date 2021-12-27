import React, { useState } from 'react';
import { Tooltip, message } from 'antd';
import { UserModelState } from '@/models/user';
import { MeetingModelState } from '@/models/meeting';
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

interface IControlBarProps {
  currentUser: UserModelState;
  meeting: MeetingModelState;
  openUsersDrawer: () => void;
  changeMicState: (s: boolean) => void;
  changeCameraState: (s: boolean) => void;
  changeShareState: (s: boolean) => void;
  leaveMeeting: () => void;
  recordMeeting: () => void;
}

const ControlBar: React.FC<IControlBarProps> = ({
  currentUser,
  meeting,
  openUsersDrawer,
  changeMicState,
  changeCameraState,
  changeShareState,
  leaveMeeting,
  recordMeeting,
}) => {
  const commonProps = {
    width: 36,
    height: 36,
    style: { background: 'transparent', margin: 0 }
  };

  const [ visible, setVisible ] = useState(false);

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
            <img src={currentUser.isSharing ? shareOffIcon : shareOnIcon} />
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
      {visible && (
        <SettingsModal visible={visible} close={() => setVisible(false)} />
      )}
    </div>
  );
};

export default ControlBar;
