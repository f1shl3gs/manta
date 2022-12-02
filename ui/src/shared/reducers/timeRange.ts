import {TimeRange} from 'src/types/timeRanges'
import {pastHourTimeRange} from 'src/constants/timeRange'
import {Action, SET_TIMERANGE} from 'src/shared/actions/timeRange'

export type TimeRangeState = TimeRange

const initialState = (): TimeRangeState => pastHourTimeRange

export const timeRangeReducer = (
  state = initialState(),
  action: Action
): TimeRangeState => {
  switch (action.type) {
    case SET_TIMERANGE:
      return action.timeRange
    default:
      return state
  }
}
