// Libraries
import React from 'react'

// Components
import {Config, Table} from '@influxdata/giraffe'

// Hooks
import {useVisXDomainSettings} from './useVisXDomainSettings'

// Types
import {getFormatter} from 'utils/vis'
import {useLineView} from 'shared/useViewProperties'

interface Props {
  table: Table
  groupKeyUnion: string[]
  children: (config: Config) => JSX.Element
}

const XYPlot: React.FC<Props> = props => {
  const {children, table, groupKeyUnion} = props
  const {xDomain, onSetXDomain, onResetXDomain} = useVisXDomainSettings()

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
    xDomain,
    onSetXDomain: onSetXDomain,
    onResetXDomain: onResetXDomain,
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
