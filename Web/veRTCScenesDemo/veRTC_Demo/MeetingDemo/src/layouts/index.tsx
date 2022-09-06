import React from 'react';
import styles from './index.less';
import { Layout } from 'antd';
import Logo from '@/components/Logo';
import VERTC from '@volcengine/rtc';

const { Header, Content } = Layout;

const BasicLayout: React.FC = (props) => {
  return (
    <Layout className={styles.layout}>
      <Header className={styles.title}>
        <Logo />
        <span className={styles.version}>
          Demo版本 V1.0.1 / SDK版本 v{VERTC.getSdkVersion()}
        </span>
      </Header>
      <Content className={styles.content}>{props.children}</Content>
    </Layout>
  );
};

export default BasicLayout;
