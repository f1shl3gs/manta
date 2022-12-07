import {TimeRange} from 'src/types/timeRanges'

export const SET_TIMERANGE = 'SET_TIMERANGE'

export const setTimeRange = (timeRange: TimeRange) =>
  ({
    type: SET_TIMERANGE,
    timeRange,
  } as const)

export type Action = ReturnType<typeof setTimeRange>
