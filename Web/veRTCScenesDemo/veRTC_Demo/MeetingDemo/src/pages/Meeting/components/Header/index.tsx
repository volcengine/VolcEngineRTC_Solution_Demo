import React, { useState, useEffect } from 'react';
import Logo from '@/components/Logo';
import { MeetingModelState, ViewMode } from '@/models/meeting';
import { UpOutlined, DownOutlined } from '@ant-design/icons';
import Utils from '@/utils/utils';
import { Badge } from 'antd';
import styles from './index.less';

interface IHeaderProps {
  username: string;
  meetingCreateAt: number;
  now: number;
  roomId: string;
  meeting: MeetingModelState;
  changeSpeakCollapse: () => void;
}

let timer: NodeJS.Timeout;

const Header: React.FC<IHeaderProps> = ({
  username,
  meetingCreateAt,
  now,
  roomId,
  meeting,
  changeSpeakCollapse,
}) => {
  const [duration, updateDuration] = useState(0);

  const startInterval = () => {
    timer = setTimeout(() => {
      updateDuration((d) => d + 1);
      startInterval();
    }, 1000);
  };

  useEffect(() => {
    return () => clearTimeout(timer);
  }, []);

  useEffect(() => {
    clearTimeout(timer);
    const d = Math.ceil((now - meetingCreateAt) / 1e9);
    updateDuration(d);
    startInterval();
  }, [meetingCreateAt, now]);

  return (
    <div className={styles.container}>
      <Logo />

      <div className={styles.center}>
        <span>房间ID: {roomId}</span>
        <span className={styles.split} />
        <span>{Utils.formatTime(duration)}</span>
        <>
          {meeting.viewMode === ViewMode.SpeakerView ? (
            <span className={styles.collapse} onClick={changeSpeakCollapse}>
              {meeting.speakCollapse ? (
                <>
                  <DownOutlined />
                  <em>展开列表</em>
                </>
              ) : (
                <>
                  <UpOutlined />
                  <em>隐藏列表</em>
                </>
              )}
            </span>
          ) : null}
        </>
        <>
          {meeting.meetingInfo.record ? (
            <span style={{ marginLeft: 100 }}>
              <Badge color="red" />
              REC
            </span>
          ) : null}
        </>
      </div>

      <div className={styles.right}>
        <span className={styles.avatar}>{username[0]}</span>
        &nbsp;
        {username}
      </div>
    </div>
  );
};

export default Header;
