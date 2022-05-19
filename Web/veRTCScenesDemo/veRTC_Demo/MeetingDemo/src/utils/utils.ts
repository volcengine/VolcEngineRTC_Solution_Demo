import { v4 as uuid } from 'uuid';
import Logger from '@/utils/Logger';
import { VerifyLoginRes } from '@/lib/socket-interfaces';
import { DeviceInstance, DeviceItems } from '@/app-interfaces';
import { throttle } from 'lodash';

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

  // static sortDevice(device: DeviceInstance[]): DeviceItems {
  //   return device?.reduce(
  //     (prev: DeviceItems, item) => {
  //       const { kind } = item;
  //       if (!prev[kind].find((i) => i.groupId === item.groupId)) {
  //         prev[kind].push(item);
  //       }
  //       return prev;
  //     },
  //     {
  //       audioinput: [],
  //       videoinput: [],
  //       audiooutput: [],
  //     }
  //   );
  // }

  static getThousand(value: number | undefined): number {
    return parseInt(((value || 0) / 1000).toString());
  }
  static protocolCheck(): void {
    if (process.env.NODE_ENV === 'development') {
      return;
    }
    const targetProtocol = 'https:';
    if (window.location.protocol !== targetProtocol) {
      window.location.href =
        targetProtocol +
        window.location.href.substring(window.location.protocol.length);
    }
  }

  /**
   * @brief 比较两个对象, 返回键值不同的键
   * @function diff
   * @param O1 对象1
   * @param O2 对象2
   * @param isolation 如果键存在于 isolation 则返回键的父键
   * @returns {[key: string]: boolean}
   */
  static diff(
    O1: any,
    O2: any,
    isolation: any = {},
    parent?: any
  ): { [key: string]: boolean } {
    let diffRes: { [key: string]: boolean } = {};
    let isolationObj: { [key: string]: boolean } = {};
    if (Array.isArray(isolation)) {
      isolation.forEach((item: string) => {
        isolationObj[item] = true;
      });
    } else isolationObj = { ...isolation };
    for (const key in O1) {
      // 如果 O2[key] 存在
      if (O2[key] !== undefined) {
        if (typeof O1[key] !== 'object') {
          // 比较
          if (O1[key] !== O2[key]) {
            // 判断是否返回父键
            if (isolationObj[key] !== undefined) diffRes[parent] = true;
            else diffRes[key] = true;
          }
        } else {
          const childDiff = Utils.diff(O1[key], O2[key], isolationObj, key);
          diffRes = { ...childDiff, ...diffRes };
        }
      }
    }
    return diffRes;
  }
  static checkMediaState(): void {
    document.body.addEventListener(
      'mousemove',
      throttle(() => {
        const remoteVideoContainers = document.querySelectorAll(
          '.remote_player_container'
        );

        Array.from(remoteVideoContainers).forEach((c) => {
          const video = c.querySelector('video');
          if (video && video.muted) {
            video.muted = false;
            video.removeAttribute('muted');
          }
          if (video && video.paused) {
            video.play();
          }
        });
      }, 1000)
    );
  }
}

export default Utils;
