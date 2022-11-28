import {TimeRange} from 'src/types/TimeRanges'

export type Action = SetTimeRange

interface SetTimeRange {
  type: 'SetTimeRange'
  payload: TimeRange
}
