import * as React from 'react';
import { v4 as uuid } from 'uuid';
import { Stream, FitType } from '@/app-interfaces';

interface IPlayerProps {
  stream: Stream;
  remote?: boolean;
}

const Player: React.FC<IPlayerProps> = (props) => {
  const id = uuid();
  React.useEffect(() => {
    if (props.stream) {
      props.stream.play(id, {
        fit: props.stream.stream.screen
          ? FitType['contain']
          : FitType['cover'],
      });
    }
  }, [props.stream.stream]);
  return (<div id={id} className={props.remote ? 'remote_player_container' : ''} style={{width: '100%', height: '100%'}}></div>);
};

export default Player;
