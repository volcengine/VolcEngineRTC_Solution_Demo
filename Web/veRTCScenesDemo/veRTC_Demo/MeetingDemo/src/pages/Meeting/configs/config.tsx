import { AppState } from '@/app-interfaces';
import { connect, bindActionCreators } from 'dva';
import { Dispatch } from '@@/plugin-dva/connect';
import { userActions } from '@/models/user';
import { meetingActions } from '@/models/meeting';
import { ConnectedProps } from 'react-redux';
import { injectIntl } from 'umi';
import React from 'react';

/**
 * @brief 获取 Redux 中信息
 * @function mapStateToProps
 * @returns
 */
const mapStateToProps = (state: AppState) => {
  return {
    currentUser: state.user,
    meeting: state.meeting,
    mc: state.meetingControl.sdk,
    settings: state.meetingSettings,
    rtc: state.rtcClientControl.rtc,
  };
};

/**
 * @brief 合并 Actions
 * @function mapDispatchToProps
 * @returns
 */
const mapDispatchToProps = (dispatch: Dispatch) => {
  return {
    dispatch,
    ...bindActionCreators({ ...userActions, ...meetingActions }, dispatch),
  };
};

const connector = connect(mapStateToProps, mapDispatchToProps);

/**
 * @brief 将 Redux 的属性注入到props
 * @function injectProps
 */
const injectProps = (comp: React.FC<any> | React.ComponentType<any>) => {
  return connector(injectIntl(comp));
};

export {
  mapStateToProps,
  mapDispatchToProps,
  connect,
  connector,
  injectProps,
  ConnectedProps,
};
