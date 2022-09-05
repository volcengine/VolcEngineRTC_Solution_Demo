export default {
  define: {
    'process.env': {
      ENV: 'test',
      SOCKETURL: 'wss://rtcio.bytedance.com',
      ICEURL:
        'https://rtc-access.bytedance.com/dispatch/v1/AccessInfo?Action=GetAccessInfo',
    },
  },
};
