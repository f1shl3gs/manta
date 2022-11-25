// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Config, Plot} from '@influxdata/giraffe'

// Hooks
import {useVisXDomainSettings} from 'src/visualization/Line/useVisXDomainSettings'
import useLineView from './useLineView'

// Utils
import {getFormatter} from 'src/shared/utils/vis'

// Types
import {XYViewProperties} from 'src/types/Dashboard'
import {VisualizationProps} from 'src/visualization'

interface Props extends VisualizationProps {
  properties: XYViewProperties
}

const Line: FunctionComponent<Props> = ({result}) => {
  const {table, fluxGroupKeyUnion} = result
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
    base: yAxisBase,
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
        fill: fluxGroupKeyUnion,
      },
    ],
  }

  return <Plot config={config} />
}

export default Line
