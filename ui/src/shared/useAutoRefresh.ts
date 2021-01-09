import constate from 'constate';
import { useEffect, useState } from 'react';
import { AutoRefresh, AutoRefreshStatus } from 'types/AutoRefresh';
import { useTimeRange } from './useTimeRange';
import { TimeRange } from '../types/TimeRanges';
import moment from 'moment';

const calculateRange = (timeRange: TimeRange) => {
  switch (timeRange.type) {
    case 'selectable-duration':
      const end = moment().unix();
      const start = end - timeRange.seconds;

      return {
        start,
        end,
        step: 14
      };

    default:
      const now = moment().unix();

      return {
        start: now - 3600,
        end: now,
        step: 14
      };
  }
};

const [AutoRefreshProvider, useAutoRefresh] = constate(
  () => {
    const { timeRange } = useTimeRange();
    const [autoRefresh, setAutoRefresh] = useState<AutoRefresh>({
      status: AutoRefreshStatus.Active,
      interval: 15
    });
    const [state, setState] = useState(() => calculateRange(timeRange));

    useEffect(() => {
      if (autoRefresh.status !== AutoRefreshStatus.Active) {
        return;
      }

      const timer = setInterval(() => {
        setState(prev => calculateRange(timeRange));
      }, autoRefresh.interval * 1000);

      return () => {
        clearInterval(timer);
      };
    }, [autoRefresh, timeRange]);

    return {
      autoRefresh,
      setAutoRefresh,
      ...state
    };
  },
  value => value
);

export {
  AutoRefreshProvider,
  useAutoRefresh
};