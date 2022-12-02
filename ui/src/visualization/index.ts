import {FromFluxResult} from '@influxdata/giraffe'
import {TimeRange} from 'src/types/timeRanges'

export interface VisualizationProps {
  result: FromFluxResult

  cellID?: string
  timeRange?: TimeRange
}
