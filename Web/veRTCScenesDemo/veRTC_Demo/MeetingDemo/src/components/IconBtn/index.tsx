import React, {ReactEventHandler, CSSProperties} from 'react';
import styles from './icon-btn.less';
export interface IconBtnProps {
  onClick?: ReactEventHandler,
  onMouseEnter?: ReactEventHandler,
  width?: number,
  height?: number,
  shape?: 'circle' | 'square',
  radius?: number,
  style?: CSSProperties
}

const IconBtn: React.FC<IconBtnProps> = props => {
  const {radius = 0, shape = 'circle', style = {}} = props;
  const finalRadius = shape === 'circle' ? '50%' : radius;
  return (
    <div
      className={styles['icon-btn']}
      onClick={props.onClick}
      onMouseEnter={props.onMouseEnter}
      style={{
        ...style,
        width: props.width || 24,
        height: props.height || 24,
        borderRadius: finalRadius
     }}
    >
      {props.children}
    </div>
  );
};

export default IconBtn;
