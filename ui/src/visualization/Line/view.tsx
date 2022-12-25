// Libraries
import React, {FunctionComponent, useMemo} from 'react'

// Components
import {Config, Plot} from '@influxdata/giraffe'

// Utils
import {getFormatter} from 'src/shared/utils/vis'

// Types
import {XYViewProperties} from 'src/types/cells'
import {VisualizationProps} from 'src/visualization'
import {DEFAULT_LINE_COLORS} from 'src/shared/constants/graphColorPalettes'
import {LineHoverDimension} from '@influxdata/giraffe/dist/types'

interface Props extends VisualizationProps {
  properties: XYViewProperties
}

const Line: FunctionComponent<Props> = ({properties, result}) => {
  const {table, fluxGroupKeyUnion} = result
  const {
    timeFormat,
    xColumn,
    yColumn,
    colors = [],
    axes: {
      x: {label: xAxisLabel},
      y: {
        label: yAxisLabel,
        base: yAxisBase,
        prefix: yAxisPrefix,
        suffix: yAxisSuffix,
      },
    },
  } = properties

  // TODO: fix table.getColumnType(xColumn)
  const xFormatter = getFormatter('time', {
    prefix: properties.axes.x.prefix,
    suffix: properties.axes.x.suffix,
    base: properties.axes.x.base,
    timeZone: 'Local',
    timeFormat,
  })

  const yFormatter = getFormatter(table.getColumnType(yColumn), {
    prefix: yAxisPrefix,
    suffix: yAxisSuffix,
    base: yAxisBase,
    timeZone: 'Local',
    timeFormat: timeFormat,
  })

  const colorHexes = useMemo(() => {
    const _colors = colors.filter(c => c.type === 'scale')
    if (_colors && _colors.length) {
      return _colors.map(color => color.hex)
    }
    return DEFAULT_LINE_COLORS.map(color => color.hex)
  }, [colors])

  const config: Config = {
    table,
    xAxisLabel,
    yAxisLabel,
    valueFormatters: {
      [xColumn]: xFormatter,
      [yColumn]: yFormatter,
    },
    layers: [
      {
        type: 'line',
        x: xColumn,
        y: yColumn,
        fill: fluxGroupKeyUnion,
        position: properties.position,
        colors: colorHexes,
        shadeBelow: !!properties.shadeBelow,
        shadeBelowOpacity: 0.08,
        hoverDimension: properties.hoverDimension as LineHoverDimension,
      },
    ],
  }

  return <Plot config={config} />
}

export default Line
