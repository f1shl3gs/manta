import constate from 'constate';
import { useState } from 'react';
import { TimeRange } from '../types/TimeRanges';
import { pastHourTimeRange } from '../constants/timeRange';

const [TimeRangeProvider, useTimeRange] = constate(
  () => {
    const [timeRange, setTimeRange] = useState<TimeRange>(pastHourTimeRange);

    return {
      timeRange,
      setTimeRange
    };
  },
  value => value
);

export {
  TimeRangeProvider,
  useTimeRange
};