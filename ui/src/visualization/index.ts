import {FromFluxResult, Table} from '@influxdata/giraffe'
import {TimeRange} from 'src/types/TimeRanges'

export interface VisualizationProps {
  result: FromFluxResult

  cellID?: string
  timeRange?: TimeRange
}
