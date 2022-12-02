import {TimeRange} from 'src/types/timeRanges'

export const SET_TIMERANGE = 'SET_TIMERANGE'

export type Action = ReturnType<typeof setTimeRange>

export const setTimeRange = (timeRange: TimeRange) => ({
  type: SET_TIMERANGE,
  timeRange,
})
