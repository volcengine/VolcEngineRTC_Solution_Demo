import React, {
  FC,
  useState,
  useCallback,
  useEffect,
} from 'react';
import { Dispatch } from '@@/plugin-dva/connect';
import { injectIntl } from 'umi';
import { connect, bindActionCreators } from 'dva';
import { userActions } from '@/models/user';
import { meetingActions } from '@/models/meeting';
import { AppState, Stream, IStreamStats } from '@/app-interfaces';
import { ConnectedProps } from 'react-redux';
import { WrappedComponentProps } from 'react-intl';
import { useInterval } from '@/utils/hook';
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
  cameraStream: Stream | null;
  remoteStreams: { [id: string]: Stream };
};

export type StatsProps = ConnectedProps<typeof connector> &
  WrappedComponentProps &
  IProps;

const StreamStats: FC<StatsProps> = ({
  remoteStreams, cameraStream,
  settings,
}) => {
  const [interval, setInterval] = useState<number | null>(null);
  const [localStats, setLocalStats] = useState<IStreamStats>();
  const [remoteStats, setRemoteStats] = useState<IStreamStats>();

  const getStreamStats = useCallback(() => {
    if (cameraStream) {
      setLocalStats({
        ...cameraStream?.getStats(),
      });
    }
    if (remoteStreams && Object.keys(remoteStreams).length) {
      const remote = Object.values(remoteStreams)?.[0];
      setRemoteStats({
        ...remote?.getStats(),
      });
    }
  }, [cameraStream, remoteStreams]);

  useInterval(
    () => {
      getStreamStats();
    },
    interval,
    { immediate: true }
  );

  useEffect(() => {
    setInterval(settings?.realtimeParam ? 1000 : null);
  }, [settings?.realtimeParam]);

  return settings?.realtimeParam ? (
    <div className={styles['status']}>
      <div>[LOCAL]</div>
      <div>
        RES：
        {`${localStats?.videoSentResolutionWidth || 0} * ${
          localStats?.videoSentResolutionHeight || 0
        }`}
      </div>
      <div>FPS：{localStats?.videoSentFrameRate || 0}</div>
      <div>
        BIT(VIDEO)：{Utils.getThousand(localStats?.videoSentBitrate)}kbps
      </div>
      <div>
        BIT(AUDIO)：{Utils.getThousand(localStats?.audioSentBitrate)}kbps
      </div>

      <>
        <div style={{ marginTop: 10 }}>[REMOTE]</div>
        <div>RTT(VIDEO)：{remoteStats?.accessDelay ?? 0}ms</div>
        <div>RTT(AUDIO)：{remoteStats?.accessDelay ?? 0}ms</div>
        <div>CPU：0%|0%</div>
        <div>
          BIT(VIDEO)：{Utils.getThousand(remoteStats?.videoSentBitrate)}kbps
        </div>
        <div>
          BIT(AUDIO): {Utils.getThousand(remoteStats?.audioSentBitrate)}kbps
        </div>
        <div>
          RES：
          {`${remoteStats?.videoReceivedResolutionWidth || 0} * ${
            remoteStats?.videoReceivedResolutionHeight || 0
          }`}
        </div>
        <div>FPS：{remoteStats?.videoSentFrameRate || 0}</div>
        <div>LOSS（VEDIO）：{remoteStats?.videoReceivedPacketsLost || 0}%</div>
        <div>LOSS(AUDIO)：{remoteStats?.audioReceivedPacketsLost || 0}%</div>
      </>
    </div>
  ) : null;
};

export default connector(injectIntl(StreamStats));