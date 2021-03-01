// Libraries
import React from 'react'

// Components
import {Config, Table} from '@influxdata/giraffe'

// Types
import {getFormatter} from 'utils/vis'
import {useLineView} from 'shared/useViewProperties'

interface Props {
  children: (config: Config) => JSX.Element
  // timeRange?: TimeRange | null
  table: Table
  groupKeyUnion: string[]
}

const XYPlot: React.FC<Props> = props => {
  const {children, table, groupKeyUnion} = props

  const {
    timeFormat,
    xColumn,
    yColumn,
    hoverDimension,
    axes: {
      x: {label: xAxisLabel, prefix: xTickPrefix, suffix: xTickSuffix},
      y: {
        label: yAxisLabel,
        base: yAxisBase,
        prefix: yAxisPrefix,
        suffix: yAxisSuffix,
      },
    },
  } = useLineView()

  const xFormatter = getFormatter('time', {
    prefix: xTickPrefix,
    suffix: xTickSuffix,
    base: '10',
    timeZone: 'Local',
    timeFormat: timeFormat,
  })

  const yFormatter = getFormatter('number', {
    prefix: yAxisPrefix,
    suffix: yAxisSuffix,
    base: yAxisBase || '10',
    timeZone: 'Local',
    timeFormat: timeFormat,
  })

  const config: Config = {
    table,
    xAxisLabel,
    yAxisLabel,
    // @ts-ignore
    valueFormatters: {
      [xColumn]: xFormatter,
      [yColumn]: yFormatter,
    },
    layers: [
      {
        hoverDimension,
        type: 'line',
        x: xColumn,
        y: yColumn,
        fill: groupKeyUnion,
      },
    ],
  }

  return children(config)
}

export default XYPlot
