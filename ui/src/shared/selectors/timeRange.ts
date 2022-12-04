import {AppState} from 'src/types/stores'
import {TimeRange} from 'src/types/timeRanges'

export const getTimeRange = (state: AppState): TimeRange => state.timeRange
