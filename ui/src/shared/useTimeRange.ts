// Libraries
import { useState } from 'react';
import constate from 'constate';

// Types
import { TimeRange } from '../types/TimeRanges';

// Constants
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