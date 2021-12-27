export const DEFAULTCONFIG = {
    resolution: {
      width: 640,
      height: 360
    },
    resolutionsText: '640x360',
    frameRate: {
      min: 15,
      max: 15
    },
    bitrate: {
      min: 800,
      max: 800
    }
  };
  
  export const RESOLUTIOIN_LIST = [
    {
      text: '160 * 160',
      val: {
        width: 160,
        height: 160
      },
      bitrateRange: {
        min: 40,
        max: 150
      }
    }, {
      text: '320 * 180',
      val: {
        width: 320,
        height: 180
      },
      bitrateRange: {
        min: 80,
        max: 350
      }
    }, {
      text: '320 * 240',
      val: {
        width: 320,
        height: 240
      },
      bitrateRange: {
        min: 100,
        max: 400
      }
    }, {
      text: '640 * 360',
      val: {
        width: 640,
        height: 360
      },
      bitrateRange: {
        min: 200,
        max: 1000
      }
    }, {
      text: '480 * 480',
      val: {
        width: 480,
        height: 480
      },
      bitrateRange: {
        min: 200,
        max: 1000
      }
    }, {
      text: '640 * 480',
      val: {
        width: 640,
        height: 480
      },
      bitrateRange: {
        min: 250,
        max: 1000
      }
    }, {
      text: '960 * 540',
      val: {
        width: 960,
        height: 540
      },
      bitrateRange: {
        min: 400,
        max: 1600
      }
    }, {
      text: '1280 * 720',
      val: {
        width: 1280,
        height: 720
      },
      bitrateRange: {
        min: 500,
        max: 2000
      }
    },
    {
      text: '1920 * 1080',
      val: {
        width: 1920,
        height: 1080
      },
      bitrateRange: {
        min: 800,
        max: 3000
      }
    }
  ];
  
  export const BITRATEMAP: { [key: string]: number[] } = {
    '160 * 160': [40, 150],
    '320 * 180': [80, 350],
    '320 * 240': [100, 400],
    '640 * 360': [200, 1000],
    '480 * 480': [200, 1000],
    '640 * 480': [250, 1000],
    '960 * 540': [400, 1600],
    '1280 * 720': [500, 2000],
    '1920 * 1080': [800, 3000]
  };
  
  export const FRAMERATE = [15, 20, 24];
  export const FEEDBACKINFO = {
    video: '视频故障',
    context: '共享内容故障',
    audio: '音频故障',
    accident: '意外结束'
  };
  
  export const SOCKETURL = process.env.SOCKETURL as string;
  export const SOCKETPATH = '/vc_control';
  export const TOASTS = {
    token_error: '服务端Token生成失败，请重试',
    mic_right: '麦克风权限已关闭，请至设备设置页开启',
    mic_setting_right: '麦克风打开失败，请检查设备',
    car_right: '摄像头权限已关闭，请至设备设置页开启',
    car_setting_right: '摄像头打开失败，请检查设备',
    screen_error: '屏幕共享失败',
    screen_not_allow: '没有屏幕共享权限',
    tick: '相同ID用户已登录，您已被强制下线！',
    mute: '你已被主持人静音',
    unmute: '主持人邀请你打开麦克风',
    give_host_error: '移交主持人失败，请重试',
    mute_error: '静音失败，请重试',
    network_error: '网络链接已断开，请检查设置',
    record: '如需录制会议，请提醒主持人开启录制',
    lock_error_track: '流中断，请刷新页面后恢复',
  };
  