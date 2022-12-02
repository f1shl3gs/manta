import {FromFluxResult} from '@influxdata/giraffe'
import {ViewProperties} from 'src/types/cells'
import {TimeRange} from 'src/types/timeRanges'

export interface VisualizationProps {
  result: FromFluxResult

  cellID?: string
  timeRange?: TimeRange
}

export interface VisualizationOptionProps {
  properties: ViewProperties
  update: (obj: any) => void
}
