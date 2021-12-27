import { v4 as uuid } from 'uuid';
import Logger from '@/utils/Logger';
import { VerifyLoginRes } from '@/lib/socket-interfaces';
import { DeviceInstance, DeviceItems } from '@/app-interfaces';

const logger = new Logger('Utils');

class Utils {
  static getQueryString(name: string): string | null {
    const reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
    const reg_rewrite = new RegExp('(^|/)' + name + '/([^/]*)(/|$)', 'i');
    const r = window.location.search.substr(1).match(reg);
    const q = window.location.pathname.substr(1).match(reg_rewrite);
    if (r != null) {
      return unescape(r[2]);
    } else if (q != null) {
      return unescape(q[2]);
    } else {
      return null;
    }
  }

  static getDeviceId(): string {
    let deviceId = Utils.getQueryString('deviceId');
    if (deviceId) {
      return deviceId;
    }
    deviceId = localStorage.getItem('deviceId');
    if (!deviceId) {
      deviceId = uuid();
      localStorage.setItem('deviceId', deviceId);
    }
    return deviceId;
  }

  /**
   * format time to {mm:ss}
   * 将时间格式化成 {mm:ss}
   * @param time number(sec)
   * @returns string 10:24
   */
  static formatTime(time: number): string {
    if (!time) {
      time = 0;
    }
    let sec: number | string = time % 60;
    if (sec < 10) {
      sec = '0' + sec;
    }
    const min = Math.floor(time / 60);
    return min + ':' + sec;
  }

  static getLoginToken(): string {
    const loginInfo = Utils.getLoginInfo();
    return loginInfo ? loginInfo.login_token : '';
  }

  static getLoginUserId(): string {
    const loginInfo = Utils.getLoginInfo();
    return loginInfo ? loginInfo.user_id : '';
  }

  static getLoginUserName(): string {
    const loginInfo = Utils.getLoginInfo();
    return loginInfo ? loginInfo.user_name : '';
  }

  static setLoginInfo(info: VerifyLoginRes): void {
    localStorage.setItem('loginInfo', JSON.stringify(info));
  }

  static removeLoginInfo(): void {
    localStorage.removeItem('loginInfo');
  }

  static getLoginInfo(): VerifyLoginRes | null {
    try {
      const loginInfo = JSON.parse(
        localStorage.getItem('loginInfo') || '{}'
      ) as VerifyLoginRes;
      return loginInfo;
    } catch (error) {
      logger.error('Utils getLoginToken');
      return null;
    }
  }

  static sortDevice(device: DeviceInstance[]): DeviceItems {
    return device?.reduce(
      (prev: DeviceItems, item) => {
        const { kind } = item;
        if (!prev[kind].find((i) => i.groupId === item.groupId)) {
          prev[kind].push(item);
        }
        return prev;
      },
      {
        audioinput: [],
        videoinput: [],
        audiooutput: [],
      }
    );
  }

  static getThousand(value: number | undefined): number {
    return parseInt(((value || 0) / 1000).toString());
  }
  static protocolCheck(): void {
    if (process.env.NODE_ENV === 'development') {
      return;
    }
    const targetProtocol = 'https:';
    if(window.location.protocol !== targetProtocol){
      window.location.href = targetProtocol + window.location.href.substring(window.location.protocol.length);
    }
  }
}

export default Utils;
