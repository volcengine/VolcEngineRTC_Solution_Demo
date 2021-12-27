import * as React from 'react';
import styles from './index.less';
// import { IRemoteAudioLevel } from '@/app-interfaces';
import { MeetingModelState } from '@/models/meeting';

interface ISpeakerViewProps {
  views: React.ReactNode[];
  screenView: React.ReactNode;
  meeting: MeetingModelState;
}

const SpeakerView: React.FC<ISpeakerViewProps> = ({
  views,
  screenView,
  meeting
}) => {
  return (
    <div className={styles.container}>
      <div
        className={
          meeting?.speakCollapse ? styles.usersViewCollapse : styles.usersView
        }
      >
        {views.map((view) => (
          <div
            className={styles.speakView}
            key={(view as React.ReactElement)?.key}
          >
            {view}
          </div>
        ))}
      </div>
      <div className={styles.screenView}>{screenView}</div>
    </div>
  );
};

export default SpeakerView;
