import * as React from 'react';
import logoImg from '/assets/images/logo.png';
import styles from './index.less';

const Logo: React.FC = () => {
  return (
    <div className={styles.container}>
      <img
        src={logoImg}
        alt="logo"
        draggable="false"
        width={242}
      />
    </div>
  );
};

export default Logo;
