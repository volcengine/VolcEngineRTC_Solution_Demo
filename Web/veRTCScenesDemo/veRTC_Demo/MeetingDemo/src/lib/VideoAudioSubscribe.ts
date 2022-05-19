import { MeetingProps, IVolume } from '@/app-interfaces';
import { SubscribeMediaType, StreamIndex, RTCStream } from '@volcengine/rtc';

type IStreams = {
  [key: string]: RTCStream;
};

type SubscribePatch = {
  subscribed: string[];
  unsubscribed: string[];
};

class VideoAudioSubscribe {
  private streamRecord: SubscribePatch = { subscribed: [], unsubscribed: [] };
  private count = 8; //订阅远端流数量 8 + 1（本地流） = 9
  private props: MeetingProps;

  private _streams: IStreams = {};

  constructor(props: MeetingProps) {
    this.props = props;
  }

  public cleanStreamRecord(): void {
    this.streamRecord = { subscribed: [], unsubscribed: [] };
  }
  public changSubscribe(newUserSort: IVolume[]): void {
    newUserSort.forEach((i, index) => {
      if (!this._streams[i.userId] || !i.userId) {
        return;
      }
      const subscribe_index = this.streamRecord?.subscribed?.findIndex(
        (item) => item === i.userId
      );
      const unsubscribe_index = this.streamRecord?.unsubscribed?.findIndex(
        (item) => item === i.userId
      );

      if (index <= this.count) {
        if (subscribe_index === -1) {
          this.streamRecord?.subscribed?.push(i.userId);
          this.subscribe(
            i.userId,
            SubscribeMediaType.AUDIO_AND_VIDEO,
            StreamIndex.STREAM_INDEX_MAIN
          );
        }
        if (unsubscribe_index !== -1) {
          this.streamRecord?.unsubscribed?.splice(unsubscribe_index, 1);
        }
      }
      if (index > this.count) {
        if (unsubscribe_index === -1) {
          this.streamRecord?.unsubscribed?.push(i.userId);
          this.subscribe(
            i.userId,
            SubscribeMediaType.AUDIO_ONLY,
            StreamIndex.STREAM_INDEX_MAIN
          );
        }
        if (subscribe_index !== -1) {
          this.streamRecord?.subscribed?.splice(subscribe_index, 1);
        }
      }
    });
  }
  public subscribe(
    userId: string,
    mediaType: number,
    streamType: number
  ): void {
    const { rtc } = this.props;
    // this._assertNotInRoom();
    rtc.engine.subscribeUserStream(userId, mediaType, streamType);
  }
  // private _assertNotInRoom() {
  //   if (!this.client || this.client.getConnectionState() !== 'CONNECTED') {
  //     throw new Error('Not in RTC Room');
  //   }
  // }

  set streams(v: IStreams) {
    this._streams = v;
  }

  addSubscribed(userId: string) {
    this.streamRecord?.subscribed.push(userId);
  }
}

export default VideoAudioSubscribe;
