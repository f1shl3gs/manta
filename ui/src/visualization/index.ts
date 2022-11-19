import {Table} from '@influxdata/giraffe'
import {TimeRange} from 'src/types/TimeRanges'

export interface VisualizationProps {
  table: Table

  cellID?: string
  timeRange?: TimeRange
}
