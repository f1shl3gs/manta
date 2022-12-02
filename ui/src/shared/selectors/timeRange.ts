import {TimeRange} from '@influxdata/clockface'
import {AppState} from 'src/types/stores'

export const getTimeRange = (state: AppState): TimeRange => state.timeRange
