import React, { useEffect, useState } from 'react';
import classnames from 'classnames';
import {round} from 'lodash'

import {
  ComponentSize,
  RemoteDataState,
  TechnoSpinner
} from '@influxdata/clockface';

interface Props {
  loading: RemoteDataState
}

const ViewLoadingSpinner: React.FC<Props> = ({ loading }) => {
  const [timerActive, setTimerActive] = useState(false);
  const [seconds, setSeconds] = useState(0);

  const timerElementClass = classnames('view-loading-spinner--timer', {
    visible: seconds > 5
  });

  const resetTimer = () => {
    setSeconds(0);
    setTimerActive(false);
  };

  useEffect(() => {
    if (loading === RemoteDataState.Done || RemoteDataState.Error) {
      resetTimer();
    }

    if (loading === RemoteDataState.Loading) {
      setTimerActive(true);
    }
  }, [loading]);

  useEffect(() => {
    if (!timerActive) {
      return;
    }

    const interval = setInterval(() => {
      setSeconds(seconds => seconds + 0.1);
    }, 100);

    return () => clearInterval(interval);
  }, [timerActive, seconds]);

  if (loading === RemoteDataState.Loading) {
    return (
      <div className={'view-loading-spinner'}>
        <TechnoSpinner diameterPixels={66} strokeWidth={ComponentSize.Medium} />
        <div className={timerElementClass}>{`${round(seconds, 1)}s`}</div>
      </div>
    );
  }

  return null;
};

export default ViewLoadingSpinner;