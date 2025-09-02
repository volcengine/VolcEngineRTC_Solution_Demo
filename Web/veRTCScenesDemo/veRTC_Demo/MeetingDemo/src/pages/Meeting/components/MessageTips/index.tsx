import React from 'react';
import Logger from '@/utils/Logger';
import { Modal, message } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { TOASTS } from '@/config';
const logger = new Logger('MessageTips');

const { confirm } = Modal;

export const showOpenMicConfirm = (callback: () => void): void => {
  confirm({
    title: TOASTS['unmute'],
    icon: <ExclamationCircleOutlined />,
    content: '',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      callback();
    },
    onCancel: () => {
      logger.debug('取消打开麦克风');
    },
  });
};

export const sendMutedInfo = (): void => {
  message.info(TOASTS['mute']);
};

export const sendInfo = (): void => {
  message.info('你已发送请求');
};

export const hostChangeInfo = (
  user_name: string,
  callback: () => void
): void => {
  Modal.info({
    title: `是否将主持人移交给: ${user_name}`,
    okText: '确定',
    cancelText: '取消',
    closable: true,
    onOk: () => {
      callback();
    },
  });
};

export const modalWarning = (
  text: keyof typeof TOASTS,
  callback?: (...args: any[]) => any
): void => {
  Modal.warning({
    title: TOASTS[text],
    okText: '确定',
    cancelText: '取消',
    closable: true,
    onOk: () => {
      callback && callback();
    },
  });
};

export const modalError = (
  text: keyof typeof TOASTS,
  callback?: (...args: any[]) => any
): void => {
  Modal.error({
    title: TOASTS[text],
    okText: '确定',
    cancelText: '取消',
    closable: true,
    onOk: () => {
      callback && callback();
    },
  });
};
