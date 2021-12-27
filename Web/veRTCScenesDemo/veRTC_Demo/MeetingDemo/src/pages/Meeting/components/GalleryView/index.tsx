import React, { useState, useEffect, useCallback } from 'react';
import chunk from 'lodash.chunk';
import debounce from 'lodash.debounce';

import styles from './index.less';

interface IGalleryViewProps {
  views: React.ReactNode[];
}


const getRowCount = (total: number) => {
  if (total >= 7) {
    return 3;
  }
  if (total >=3) {
    return 2;
  }
  return 1;
};

const getColCount = (total: number) => {
  if ( total >=5 ) {
    return 3;
  }
  if (total >= 2) {
    return 2;
  }
  return 1;
};

const GalleryView: React.FC<IGalleryViewProps> = ({ views }) => {
  const [layout, updateLayout] = useState({
    rows: getRowCount(views.length) || 1,
    cols: getColCount(views.length) || 1,
  });

  const updateView = useCallback( // eslint-disable-line react-hooks/exhaustive-deps
    debounce((len)=> {
      updateLayout({
        rows: getRowCount(len),
        cols: getColCount(len),
      });
    },500),
  []);

  useEffect(() => {
    const len = views.length;
    if (len) {
      updateView(len);
    }
  }, [views, updateView]);

  const renderViews = () => {
    if (views.length) {
      const groups = chunk(views, layout.cols);
      const ret = [];
      for (let i = 0; i < groups.length; i++) {
        const row = (
          <div
            key={i}
            className={styles.galleryRow}
            style={{
              height: `${100 / layout.rows}%`,
            }}
          >
            {groups[i].map((view) => (
              <div
                className={styles.galleryView}
                key={(view as React.ReactElement)?.key}
                style={{
                  padding: groups.length ===1?'0px':'12px',
                  width: `${100 / layout.cols}%`,
                }}
              >
                {view}
              </div>
            ))}
          </div>
        );
        ret.push(row);
      }
      return ret;
    }
    return [];
  };

  return <div className={styles.container}>{renderViews()}</div>;
};

export default GalleryView;

