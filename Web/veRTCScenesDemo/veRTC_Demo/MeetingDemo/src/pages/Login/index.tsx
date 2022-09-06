import React, { Component, ReactNode } from 'react';
import { Dispatch } from '@@/plugin-dva/connect';
import { injectIntl, history } from 'umi';
import { connect, bindActionCreators } from 'dva';
import { ConnectedProps } from 'react-redux';
import { WrappedComponentProps } from 'react-intl';
import { userActions } from '@/models/user';
import { AppState } from '@/app-interfaces';
import styles from './index.less';
import { Button, Form, Input, Checkbox, message } from 'antd';
import { FormInstance } from 'antd/es/form';
import Logger from '@/utils/Logger';
import Utils from '@/utils/utils';
import { StoreValue } from 'rc-field-form/lib/interface';
import titleLogo from '/assets/images/titleLogo.png';
import type { GetAppIDResponse } from '@/lib/socket-interfaces';

const logger = new Logger('login');

enum ERROR_TYPES {
  VALID,
  EMPTY_STRING,
  INVALID_CHARACTERS,
}

const messages = {
  codeErrType: {
    1: '请填写验证码',
    2: '验证码输入有误，请重新输入',
  },
  phoneErrType: {
    1: '请填写手机号',
    2: '手机号输入有误，请重新输入',
  },
};

export interface LoginState {
  time: number;
  loading: boolean;
}

function mapStateToProps(state: AppState) {
  return {
    currentUser: state.user,
    mc: state.meetingControl.sdk,
    rtc: state.rtcClientControl.rtc,
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators(userActions, dispatch),
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export type LoginProps = ConnectedProps<typeof connector> &
  WrappedComponentProps;

export interface FormProps {
  cell_phone: string;
  code: string;
  agree: boolean;
}

export class Login extends Component<LoginProps, LoginState> {
  constructor(props: LoginProps) {
    super(props);
    this.state = {
      time: 60,
      loading: false,
    };
  }

  formRef = React.createRef<FormInstance>();

  TelIdpLoginSms = (): void => {
    this.props.mc?.checkSocket().then(() => {
      this.formRef.current?.validateFields(['cell_phone']).then((values) => {
        try {
          this.props.mc
            ?.getPhoneVerifyCode({
              cell_phone: values.cell_phone,
              country_code: '86',
            })
            .then(() => {
              if (this.state.time !== 0) {
                this.count();
                this.setState({ loading: true });
              }
            })
            .catch((error) => {
              message.error(`${error}`);
            });
        } catch (error) {
          message.error(`${error}`);
        }
      });
    });
  };

  onFinish = ({ cell_phone, code, agree }: FormProps): void => {
    if (!agree) {
      return;
    }
    try {
      this.props.mc?.checkSocket().then(() => {
        this.props.mc
          ?.verifyLoginSms({
            cell_phone,
            country_code: '86',
            code,
          })
          .then((res) => {
            Utils.setLoginInfo(res);
            this.props.setLogged(true);
            this.props.mc?.getAppID({}).then((app?: GetAppIDResponse) => {
              if (!app) {
                return;
              }
              this.props.setAppId(app.app_id);
              this.props.rtc.init({
                config: {
                  appId: app.app_id,
                  uid: res.user_id,
                },
              });
              history.push('/');
            });
          })
          .catch((error) => {
            message.error(`${error}`);
          });
      });
    } catch (error) {
      message.error(`${error}`);
      logger.warn('verifyLoginSms error', error);
    }
  };

  count(): void {
    let { time } = this.state;
    const siv = setInterval(() => {
      this.setState({ time: time-- }, () => {
        if (time <= -1) {
          clearInterval(siv);
          this.setState({ loading: false, time: 60 });
        }
      });
    }, 1000);
  }

  validator(
    value: StoreValue,
    errorTypeKey: 'phoneErrType' | 'codeErrType',
    regRes: boolean
  ): Promise<void | any> | void {
    let result: Promise<Error | void>;
    if (!value || regRes) {
      const _value = value
        ? ERROR_TYPES.INVALID_CHARACTERS
        : ERROR_TYPES.EMPTY_STRING;
      result = Promise.reject(new Error(messages[errorTypeKey][_value]));
    } else {
      result = Promise.resolve();
    }
    return result;
  }

  render(): ReactNode {
    const { loading, time } = this.state;

    return (
      <div className={styles.container}>
        <div className={styles.main}>
          <div className={styles['main-title']}>
            <img width={400} src={titleLogo} alt="logo" />
          </div>
          <Form ref={this.formRef} onFinish={this.onFinish}>
            <Form.Item
              name="cell_phone"
              validateTrigger="onChange"
              rules={[
                {
                  required: true,
                  validator: (_, value) => {
                    const res = !/^[1][3,4,5,7,8,9][0-9]{9}$/.test(value);
                    return this.validator(value, 'phoneErrType', res);
                  },
                },
              ]}
            >
              <Input
                autoComplete="off"
                placeholder="输入手机号"
                className={styles['login-input']}
              />
            </Form.Item>
            <Form.Item
              name="code"
              rules={[
                {
                  required: true,
                  validator: (_, value) => {
                    const res = !/^\d{6}$/.test(value);
                    return this.validator(value, 'codeErrType', res);
                  },
                },
              ]}
            >
              <Input
                autoComplete="off"
                placeholder="输入验证码"
                suffix={
                  <Button
                    style={{ width: 112 }}
                    onClick={this.TelIdpLoginSms}
                    disabled={loading}
                  >
                    {loading ? `${time}s` : '获取验证码'}
                  </Button>
                }
              />
            </Form.Item>
            <Form.Item
              name="agree"
              valuePropName="checked"
              wrapperCol={{ span: 24 }}
              className={styles['login-agree']}
              rules={[
                {
                  required: true,
                  message: '请先阅读并同意',
                },
              ]}
            >
              <Checkbox>
                已阅读并同意
                <a
                  href="https://www.volcengine.com/docs/6348/68917"
                  target="_blank"
                  rel="noreferrer"
                >
                  《服务协议》
                </a>
                和
                <a
                  href="https://www.volcengine.com/docs/6348/68918"
                  target="_blank"
                  rel="noreferrer"
                >
                  《隐私权政策》
                </a>
              </Checkbox>
            </Form.Item>
            <Form.Item
              noStyle
              shouldUpdate={(prevValues, curValues) =>
                prevValues.agree !== curValues.agree
              }
            >
              {({ getFieldValue }) => {
                return (
                  <Form.Item>
                    <Button
                      disabled={!getFieldValue('agree')}
                      htmlType="submit"
                      type="primary"
                      className={styles['login-check']}
                    >
                      登录
                    </Button>
                  </Form.Item>
                );
              }}
            </Form.Item>
          </Form>
        </div>
      </div>
    );
  }
}

export default connector(injectIntl(Login));
