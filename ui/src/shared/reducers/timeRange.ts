import {TimeRange} from 'src/types/TimeRanges'
import {pastHourTimeRange} from 'src/constants/timeRange'
import {Action} from 'src/shared/actions/timeRange'

export type TimeRangeState = TimeRange

const initialState = (): TimeRangeState => pastHourTimeRange

export const timeRangeReducer = (
  state = initialState(),
  action: Action
): TimeRangeState => {
  switch (action.type) {
    case 'SetTimeRange':
      return action.payload
    default:
      return state
  }
}
