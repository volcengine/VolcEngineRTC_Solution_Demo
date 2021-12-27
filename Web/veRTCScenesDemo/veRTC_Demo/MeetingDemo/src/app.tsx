import { message } from 'antd';
import './app.less';
export const dva = {
  config: {
    onError(e: Error) {
      message.error(e.message, 3);
    },
  },
};
