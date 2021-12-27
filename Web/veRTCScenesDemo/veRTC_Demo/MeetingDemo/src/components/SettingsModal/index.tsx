import React, { FC, useState, useEffect, useMemo } from 'react';
import { Modal, Form, Select, Switch, Slider, Row, Col, notification } from 'antd';
import { connect, bindActionCreators } from 'dva';
import { injectIntl } from 'umi';
import { ConnectedProps } from 'react-redux';
import { WrappedComponentProps } from 'react-intl';
import styles from './index.less';
import { RESOLUTIOIN_LIST, FRAMERATE, BITRATEMAP } from '@/config';
import { AppState, DeviceInstance, DeviceItems } from '@/app-interfaces';
import { Dispatch } from '@@/plugin-dva/connect';
import { meetingSettingsActions } from '@/models/meeting-settings';
import { HistoryVideoRecord } from '@/lib/socket-interfaces';
import deleteIcon from '/assets/images/deleteIcon.png';
import moment from 'moment';
import RTC from '@/sdk/VRTC.esm.min.js';
import Logger from '@/utils/Logger';
import Utils from '@/utils/utils';

const logger = new Logger('Settings');

function mapStateToProps(state: AppState) {
  return {
    mc: state.meetingControl.sdk,
    settings: state.meetingSettings
  };
}

function mapDispatchToProps(dispatch: Dispatch) {
  return {
    dispatch,
    ...bindActionCreators(meetingSettingsActions, dispatch)
  };
}

const connector = connect(mapStateToProps, mapDispatchToProps);

export type SettingsModalProps = ConnectedProps<typeof connector> & WrappedComponentProps & { visible: boolean, close: () => void };

const commonCol = {
  labelCol: { span: 8 },
  wrapperCol: { span: 10 },
};

const SettingsModal: FC<SettingsModalProps> = (props) => {
  const [form] = Form.useForm();
  const [videoList, setVideoList] = useState<HistoryVideoRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [devices, setDevices] = useState<DeviceItems>();

  const {
    visible,
    close,
    setStreamSettings,
    setScreenStreamSettings,
    setMic,
    setCamera,
    setRealtimeParam,
    mc,
    settings
  } = props;

  const initialValues = useMemo(() => {
    const {
      streamSettings: { bitrate: BPS, frameRate: FPS, resolution },
      screenStreamSettings: {
        bitrate: shareBPS,
        frameRate: shareFPS,
        resolution: shareResolution,
      },
      mic,
      camera,
      realtimeParam,
    } = settings;
    return {
      resolution: `${resolution.width} * ${resolution.height}`,
      shareResolution: `${shareResolution.width} * ${shareResolution.height}`,
      FPS: FPS.max,
      shareFPS: shareFPS.max,
      BPS: BPS.max,
      shareBPS: shareBPS.max,
      mic: mic || devices?.audioinput[0]?.deviceId,
      camera : camera || devices?.videoinput[0]?.deviceId,
      realtimeParam,
    };
  }, [settings, devices]);

  useEffect(() => {
    form.setFieldsValue({
      ...initialValues,
    });
  }, [form, initialValues]);

  const formatStreamSettings = (res: string, fps: number, bps: number, bpsMin: number) => {
    return {
      resolution: {
        width: parseInt(res.split(' * ')[0]),
        height: parseInt(res.split(' * ')[1]),
      },
      frameRate: { min: 10, max: fps },
      bitrate: {
        min: bpsMin,
        max: bps,
      },
    };
  };

  const onOk = () => {
    const data = form.getFieldsValue(true);
    setStreamSettings(formatStreamSettings(data.resolution, data.FPS, data.BPS, 250));
    setScreenStreamSettings(formatStreamSettings(data.shareResolution, data.shareFPS, data.shareBPS, 800));
    setMic(data.mic);
    setCamera(data.camera);
    setRealtimeParam(data.realtimeParam);
    close();
  };

  const getHistoryVideoRecord = () => {
    mc?.checkSocket().then(() => {
      mc?.getHistoryVideoRecord().then((res) => {
        setVideoList(res);
        setLoading(false);
      });
    });
  };

  useEffect(() => {
    RTC.getDevices(
      (devices: DeviceInstance[]) => {
        setDevices(Utils.sortDevice(devices));
      },
      (err: Error) => {
        logger.error('err', err);
      }
    );
    getHistoryVideoRecord();
  }, []);

  const myVideoList = useMemo(() => {
    return videoList?.filter(item => item.video_holder);
  }, [videoList]);

  return (
    <Modal
      title="会议设置"
      visible={visible}
      width={788}
      className={styles['settings-modal']}
      onCancel={close}
      onOk={onOk}
    >
      <Form form={form} labelCol={{ span: 4 }} initialValues={initialValues}>
        <Row>
          <Col span={12}>
            <Form.Item label="分辨率" name="resolution" {...commonCol}>
              <Select>
                {RESOLUTIOIN_LIST.map((item) => (
                  <Select.Option key={item.text} value={item.text}>
                    {item.text}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item
              label="屏幕共享分辨率"
              name="shareResolution"
              {...commonCol}
            >
              <Select>
                {RESOLUTIOIN_LIST.map((item) => (
                  <Select.Option key={item.text} value={item.text}>
                    {item.text}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
        </Row>
        <Row>
          <Col span={12}>
            <Form.Item label="帧率" name="FPS" {...commonCol}>
              <Select>
                {FRAMERATE.map((item) => (
                  <Select.Option key={item} value={item}>
                    {`${item} fps`}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item label="屏幕共享帧率" name="shareFPS" {...commonCol}>
              <Select>
                {FRAMERATE.map((item) => (
                  <Select.Option key={item} value={item}>
                    {`${item} fps`}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
        </Row>
        <Row>
          <Col span={12}>
            <Form.Item label="码率" {...commonCol} wrapperCol={{ span: 16 }}>
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <div
                  style={{
                    display: 'flex',
                    width: 160,
                    justifyContent: 'space-between',
                  }}
                >
                  <Form.Item
                    noStyle
                    shouldUpdate={(prevValues, currentValues) =>
                      prevValues.resolution !== currentValues.resolution
                    }
                  >
                    {() => {
                      const res = form.getFieldValue('resolution') as string;
                      const range = BITRATEMAP[res];
                      return (
                        <Form.Item noStyle name="BPS">
                          <Slider
                            min={range[0]}
                            max={range[1]}
                            style={{ width: 82 }}
                            tooltipVisible={false}
                          />
                        </Form.Item>
                      );
                    }}
                  </Form.Item>
                  <Form.Item
                    noStyle
                    shouldUpdate={(prevValues, currentValues) =>
                      prevValues.BPS !== currentValues.BPS
                    }
                  >
                    {() => (
                      <div className={styles['slider-number']}>
                        {form.getFieldValue('BPS')}
                      </div>
                    )}
                  </Form.Item>
                </div>
                <div style={{ marginLeft: 8 }}>kbps</div>
              </div>
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item
              label="屏幕共享码率"
              {...commonCol}
              wrapperCol={{ span: 16 }}
            >
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <div
                  style={{
                    display: 'flex',
                    width: 160,
                    justifyContent: 'space-between',
                  }}
                >
                  <Form.Item
                    noStyle
                    shouldUpdate={(prevValues, currentValues) =>
                      prevValues.shareResolution !==
                      currentValues.shareResolution
                    }
                  >
                    {() => {
                      const res = form.getFieldValue(
                        'shareResolution'
                      ) as string;
                      const range = BITRATEMAP[res];
                      return (
                        <Form.Item noStyle name="shareBPS">
                          <Slider
                            min={range[0]}
                            max={range[1]}
                            style={{ width: 82 }}
                            tooltipVisible={false}
                          />
                        </Form.Item>
                      );
                    }}
                  </Form.Item>
                  <Form.Item
                    noStyle
                    shouldUpdate={(prevValues, currentValues) =>
                      prevValues.shareBPS !== currentValues.shareBPS
                    }
                  >
                    {() => (
                      <div className={styles['slider-number']}>
                        {form.getFieldValue('shareBPS')}
                      </div>
                    )}
                  </Form.Item>
                </div>
                <div style={{ marginLeft: 8 }}>kbps</div>
              </div>
            </Form.Item>
          </Col>
        </Row>
        <Form.Item
          label="麦克风"
          name="mic"
          wrapperCol={{ span: 5 }}
        >
          <Select dropdownMatchSelectWidth={false}>
            {devices?.audioinput.map((item) => (
              <Select.Option value={item.deviceId} key={item.deviceId}>
                {item.label}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="摄像头"
          name="camera"
          wrapperCol={{ span: 5 }}
        >
          <Select dropdownMatchSelectWidth={false}>
            {devices?.videoinput.map((item) => (
              <Select.Option value={item.deviceId} key={item.deviceId}>
                {item.label}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="查看历史会议"
          name="history"
          wrapperCol={{ span: 10 }}
        >
          <Select
            placeholder="选择历史会议点击链接查看"
            loading={loading}
            onSelect={(url) => window.open(url as string, '_blank')}
          >
            {videoList?.map((item) => {
              return (
                <Select.Option key={item.created_at} value={item.download_url}>
                  {moment(item.created_at / 1000000).format(
                    'YYYY-MM-DD HH:mm:ss'
                  )}
                </Select.Option>
              );
            })}
          </Select>
        </Form.Item>
        <Form.Item label="我的云录制" name="record" wrapperCol={{ span: 10 }}>
          <Select
            placeholder="会议录制者有权在此处查看和删除录像"
            loading={loading}
            onSelect={(url) => window.open(url as string, '_blank')}
          >
            {myVideoList?.map((item) => (
              <Select.Option
                key={item.created_at}
                value={item.download_url}
                className={styles['my-video-list']}
              >
                <div className={styles['my-video-list-item']}>
                  <div>
                    {moment(item.created_at / 1000000).format(
                      'YYYY-MM-DD HH:mm:ss'
                    )}
                  </div>
                  <img
                    src={deleteIcon}
                    width={12}
                    height={12}
                    onClick={(e) => {
                      e.stopPropagation();
                      setLoading(true);
                      mc?.deleteVideoRecord({ vid: item.vid })
                        .then(() => {
                          notification.success({ message: '删除成功' });
                        })
                        .then(() => getHistoryVideoRecord())
                        .catch((err) => {
                          notification.error({
                            message: '删除录制记录失败',
                            description: err,
                          });
                          setLoading(false);
                        });
                    }}
                  />
                </div>
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="实时视频参数"
          name="realtimeParam"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default connector(injectIntl(SettingsModal));
