import React, { useLayoutEffect, useRef, useEffect } from 'react';
import { StreamIndex } from '@volcengine/rtc';

export interface VideoPlayerProps {
  userId: string;
  setRemoteVideoPlayer: any;
  stream: any;
}

export const ShareView: React.FC<{
  localCaptureSuccess: boolean;
  rtc: any;
  sharingUser: any;
  me: boolean;
}> = ({ rtc, sharingUser, me, localCaptureSuccess }) => {
  const userId = sharingUser.user_id;

  useEffect(() => {
    if (localCaptureSuccess) {
      if (me) {
        rtc.setLocalVideoPlayer(
          StreamIndex.STREAM_INDEX_SCREEN,
          'screen-player'
        );
      } else {
        rtc.setRemoteVideoPlayer(
          StreamIndex.STREAM_INDEX_SCREEN,
          userId,
          'screen-player',
          { isScreen: true }
        );
      }
    }
  }, [me, rtc, sharingUser.user_id, userId, localCaptureSuccess]);

  return (
    <div
      id="screen-player"
      style={{
        width: '100%',
        height: '100%',
        background: '#000',
      }}
    ></div>
  );
};

export const LocalPlayer: React.FC<{
  localCaptureSuccess: boolean;
  rtc: any;
  renderDom: string;
}> = ({ localCaptureSuccess, rtc, renderDom }) => {
  useEffect(() => {
    if (localCaptureSuccess) {
      rtc.setLocalVideoPlayer(StreamIndex.STREAM_INDEX_MAIN, renderDom);
    }
  }, [localCaptureSuccess, rtc, renderDom]);
  return (
    <div
      style={{
        width: '100%',
        height: '100%',
        background: '#000',
      }}
      id={renderDom}
    ></div>
  );
};

export const MediaPlayer: React.FC<VideoPlayerProps> = ({
  userId,
  setRemoteVideoPlayer,
  stream,
}) => {
  const dom = useRef<any>();
  useLayoutEffect(() => {
    if (setRemoteVideoPlayer && dom) {
      setRemoteVideoPlayer(
        StreamIndex.STREAM_INDEX_MAIN,
        userId,
        `remoteStream_${userId}`,
        stream
      );
    }
  }, [setRemoteVideoPlayer, stream, userId]);

  return (
    <div
      style={{width: '100%', height: '100%', position: 'relative'}}
      className='remoteStream'
      id={`remoteStream_${userId}`}
      ref={dom}
    >
    </div>
  );
};
