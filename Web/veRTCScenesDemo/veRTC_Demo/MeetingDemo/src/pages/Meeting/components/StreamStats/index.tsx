import React, { FC, useMemo } from 'react';
import { Dispatch } from '@@/plugin-dva/connect';
import { injectIntl } from 'umi';
import { connect, bindActionCreators } from 'dva';
import { userActions } from '@/models/user';
import { meetingActions } from '@/models/meeting';
import {
  AppState,
  IMeetingState,
  LocalStats,
  RemoteStats,
} from '@/app-interfaces';
import { ConnectedProps } from 'react-redux';
import { WrappedComponentProps } from 'react-intl';
import Utils from '@/utils/utils';
import styles from './index.less';

type StreamStats = {
  accessDelay: number;
  videoSentFrameRate: number;
  videoSentResolutionWidth: number;
  videoSentResolutionHeight: number;
  videoSentBitrate: number;
  audioSentBitrate: number;
  videoReceivedPacketsLost: number;
  audioReceivedPacketsLost: number;
  videoReceivedResolutionHeight: number;
  videoReceivedResolutionWidth: number;
};

function mapStateToProps(state: AppState) {
  return {
    user: state.user,
    settings: state.meetingSettings,
    meeting: state.meeting,
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators({ ...userActions, ...meetingActions }, dispatch),
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

type IProps = {
  streamStatses: IMeetingState['streamStatses'];
};

export type StatsProps = ConnectedProps<typeof connector> &
  WrappedComponentProps &
  IProps;

const StreamStats: FC<StatsProps> = ({ settings, streamStatses }) => {
  const local = useMemo(
    () => streamStatses.local,
    [streamStatses.local]
  ) as LocalStats;

  const remote = useMemo(
    () => streamStatses.remoteStreams,
    [streamStatses.remoteStreams]
  ) as { [key: string]: RemoteStats };

  return settings?.realtimeParam ? (
    <div className={styles['status']}>
      <div>[LOCAL]</div>
      <div>
        RES：
        {`${local?.videoStats.encodedFrameWidth || 0} * ${
          local?.videoStats.encodedFrameHeight || 0
        }`}
      </div>
      <div>FPS：{local?.videoStats.inputFrameRate || 0}</div>
      <div>
        BIT(VIDEO)：{Utils.getThousand(local?.videoStats.sentKBitrate)}kbps
      </div>
      <div>
        BIT(AUDIO)：{Utils.getThousand(local?.audioStats.sentKBitrate)}kbps
      </div>

      {Object.keys(remote).map((i) => {
        const remoteStats = remote[i];
        return (
          <>
            <div style={{ marginTop: 10 }}>[REMOTE]{i}</div>
            <div>RTT(VIDEO)：{remoteStats?.videoStats.e2eDelay ?? 0}ms</div>
            <div>RTT(AUDIO)：{remoteStats?.audioStats.e2eDelay ?? 0}ms</div>
            <div>CPU：0%|0%</div>
            <div>
              BIT(VIDEO)：
              {Utils.getThousand(remoteStats?.videoStats.receivedKBitrate)}kbps
            </div>
            <div>
              BIT(AUDIO):{' '}
              {Utils.getThousand(remoteStats?.audioStats.receivedKBitrate)}kbps
            </div>
            <div>
              RES：
              {`${remoteStats?.videoStats.width || 0} * ${
                remoteStats?.videoStats.height || 0
              }`}
            </div>
            <div>
              FPS：{remoteStats?.videoStats.decoderOutputFrameRate || 0}
            </div>
            <div>
              LOSS（VEDIO）：{remoteStats?.videoStats.videoLossRate || 0}%
            </div>
            <div>
              LOSS(AUDIO)：{remoteStats?.audioStats.audioLossRate || 0}%
            </div>
          </>
        );
      })}
    </div>
  ) : null;
};

export default connector(injectIntl(StreamStats));
