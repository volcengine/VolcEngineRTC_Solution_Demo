import React, { FC, useState } from 'react';
import { Modal, Row, Col, Input, Button } from 'antd';
import * as styles from './index.less';

const { TextArea } = Input;

const questions = ['视频卡顿', '共享内容故障', '音频卡顿', '其他问题'];

const ToggleButton = () => {
  const [selected, setSelect] = useState<string[]>([]);

  const clickDom = (v: string) => {
    const _index = selected.indexOf(v);
    if (_index === -1) {
      setSelect([...selected, v]);
    } else {
      selected.splice(_index, 1);
      setSelect([...selected]);
    }
  };

  return (
    <>
      {questions.map((item) => (
        <span
          key={item}
          onClick={() => clickDom(item)}
          className={selected.indexOf(item) == -1 ? '' : styles['span-selected']}
        >
          {item}
        </span>
      ))}
    </>
  );
};

const FeedBack: FC<{ status: string }> = ({ status }) => {
  const [ visible, setVisible ] = useState(status === 'end');
  const [detailVisible, setDetailVisible] = useState(false);

  const openDetail = () => {
    setVisible(false);
    setDetailVisible(true);
  };

  return (
    <>
      <Modal
        className={styles['feedback']}
        width={320}
        visible={visible}
        onCancel={() => setVisible(false)}
        footer={null}
        centered
      >
        <div className={styles['title']}>本次通话体验如何？</div>
        <Row align="middle">
          <Col
            span={12}
            style={{ textAlign: 'center', borderRight: '1px solid #E5E6EB' }}
          >
            <a className={styles['good']} onClick={() => setVisible(false)}></a>
          </Col>
          <Col span={12} style={{ textAlign: 'center' }}>
            <a className={styles['bad']} onClick={openDetail}></a>
          </Col>
        </Row>
      </Modal>
      <Modal
        className={styles['feedback']}
        visible={detailVisible}
        footer={null}
        onCancel={() => setDetailVisible(false)}
        centered
      >
        <div className={styles['title']}>具体问题反馈</div>
        <div className={styles['questions']}>
          <ToggleButton />
        </div>
        <div className={styles['title']}>其他问题反馈</div>
        <div className={styles['text']}>
          <TextArea
            placeholder="最多可输入500个字符"
            showCount
            maxLength={500}
            autoSize={{ minRows: 3, maxRows: 4 }}
          />
        </div>
        <div className={styles['bottom']}>
          <Button
            style={{ background: '#F2F3F8', borderColor: 'F2F3F8' }}
            onClick={() => setDetailVisible(false)}
          >
            取消
          </Button>
          <Button
            type="primary"
            style={{ marginLeft: 12 }}
            onClick={() => setDetailVisible(false)}
          >
            提交
          </Button>
        </div>
      </Modal>
    </>
  );
};

export default React.memo(FeedBack);
