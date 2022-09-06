import React, { useEffect, useState, useCallback } from 'react';
import { Tooltip } from 'antd';
import { v4 as uuid } from 'uuid';
import { useSize } from '@umijs/hooks';
import Logger from '@/utils/Logger';
import type { ActiveMeetingUser } from '@/pages/Meeting/components/MeetingViews';
import shareOffIcon from '/assets/images/shareOffIcon.png';
import styles from './index.less';

type IViewProps = ActiveMeetingUser & {
  avatarOnCamOff?: React.ReactElement;
  sharingId?: string;
  volume: number;
  player: JSX.Element | undefined;
  sharingView: boolean;
};

const logger = new Logger('view');

const View: React.FC<IViewProps> = ({
  me = false,
  speaking,
  is_host: isHost,
  is_sharing: isSharing,
  is_camera_on: isCameraOn,
  is_mic_on: isMicOn,
  user_name: username = '',
  avatarOnCamOff = null,
  sharingId = '',
  user_id,
  volume,
  player,
  sharingView,
}) => {
  const [id] = useState<string>(uuid());
  const [layoutSize, containerRef] = useSize<HTMLDivElement>();
  const [avatarSize, updateAvatarSize] = useState(24);

  const handleResize = () => {
    if (layoutSize.width && layoutSize.height) {
      const { width, height } = layoutSize;
      const shortRect = width > height ? height : width;
      const halfShortRect = Math.round(shortRect / 3);
      if (halfShortRect > 150) {
        updateAvatarSize(150);
      } else if (halfShortRect < 24) {
        updateAvatarSize(24);
      } else {
        updateAvatarSize(halfShortRect);
      }
    }
  };

  useEffect(() => {
    if (layoutSize.width && layoutSize.height) {
      handleResize();
    }
  }, [layoutSize]);

  const hasVolume = useCallback(() => {
    let res: boolean;
    if (me && isMicOn && volume > 0) {
      res = true;
    } else {
      res = volume > 0;
    }
    return res;
  }, [volume, isMicOn, me]);

  const render = () => {
    if (!sharingView && !isCameraOn) {
      return (
        <div className={styles.layoutWithoutCamera}>
          <div
            className={styles.avatar}
            style={{
              width: avatarSize,
              height: avatarSize,
              borderRadius: avatarSize / 2,
              lineHeight: `${avatarSize}px`,
              fontSize: avatarSize / 2,
              fontWeight: 600,
              border: hasVolume() && speaking ? '1px solid #3370FF' : 'none',
            }}
          >
            {avatarOnCamOff ? avatarOnCamOff : username[0]}
          </div>

          <div className={styles.username}>
            <Tooltip title={username}>
              {username.length > 10 ? `${username.slice(0, 10)}...` : username}
            </Tooltip>
            {me && <span className={styles.me}>我</span>}
            {isHost && <span className={styles.me}>主持人</span>}
          </div>
        </div>
      );
    }
  };

  return (
    <div
      className={styles.container}
      ref={containerRef}
      style={{ border: hasVolume() && speaking ? '1px solid #3370FF' : 'none' }}
    >
      {render()}

      <div
        className={styles.streamContainer}
        id={id}
        style={{
          display: sharingView || isCameraOn ? 'block' : 'none',
        }}
      >
        {player}
        <span className={styles.username2}>
          {sharingId === user_id && (
            <img src={shareOffIcon} style={{ width: 14 }} />
          )}
          &nbsp;
          {username}
          {me && <span className={styles.me}>我</span>}
          {isHost && <span className={styles.me}>主持人</span>}
        </span>
      </div>
    </div>
  );
};

export default View;
