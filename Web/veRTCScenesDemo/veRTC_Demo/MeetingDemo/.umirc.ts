import { defineConfig } from 'umi';

export default defineConfig({
  nodeModulesTransform: {
    type: 'none',
  },
  mfsu: {},
  antd: {},
  locale: {
    default: 'en-US',
    antd: true,
  },
  dva: {
    immer: true,
    hmr: false,
    skipModelValidate: true,
  },
  routes: [
    {
      path: '/login',
      component: '../pages/Login',
    },
    {
      path: '/solution/meeting/login',
      component: '../pages/Login',
    },
    {
      exact: true,
      path: '/',
      wrappers: ['@/wrappers/auth'],
      component: '../layouts/index',
      routes: [{ path: '/', component: '../pages/JoinRoom/index' }],
    },
    {
      path: '/meeting',
      wrappers: ['@/wrappers/auth'],
      component: '../pages/Meeting',
    },
  ],
  theme: {
    'input-border-color': 'transparent',
    'input-hover-border-color': 'transparent',
    'input-bg': '#F2F3F5',
  },
  fastRefresh: {},
  chainWebpack: function (memo) {
    memo.resolve.extensions.add('.tsx');
  },
});
