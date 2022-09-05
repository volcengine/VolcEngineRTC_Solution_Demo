import * as React from 'react';
import { Modal, Button } from 'antd';
import exclamationMarkIcon from '/assets/images/exclamationMarkIcon.png';
import styles from './index.less';

interface ILeavingConfirmProps {
  visible?: boolean;
  isHost?: boolean;
  endMeeting: () => void;
  leaveMeeting: () => void;
  cancel: () => void;
}

const LeavingConfirm: React.FC<ILeavingConfirmProps> = ({
  visible,
  isHost,
  endMeeting,
  leaveMeeting,
  cancel,
}) => {
  return (
    <Modal
      className={styles.container}
      visible={visible}
      centered
      width={400}
      footer={null}
      closable={false}
      bodyStyle={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
      }}
    >
      <div>
        <div className={styles.mark}>
          <img src={exclamationMarkIcon} alt="exclamationMarkIcon" />
        </div>
        <span className={styles.title}>
          {isHost
            ? '请移交主持人给指定参会者，方能离开会议'
            : '请再次确认是否离开会议？'}
        </span>
      </div>
      {isHost && (
        <Button
          className={styles.btn}
          type="primary"
          danger
          onClick={endMeeting}
        >
          结束全部会议
        </Button>
      )}
      <Button
        className={styles.btn}
        type="primary"
        danger
        disabled={isHost}
        onClick={leaveMeeting}
      >
        离开会议
      </Button>
      <Button className={styles.btn} onClick={cancel}>
        取消
      </Button>
    </Modal>
  );
};

export default LeavingConfirm;
