import React, {FunctionComponent, useMemo} from 'react'
import {VisualizationProps} from 'src/visualization'
import {LinePlusSingleStatViewProperties} from 'src/types/cells'
import {
  Config,
  getLatestValues,
  Plot,
  SINGLE_STAT_SVG_NO_USER_SELECT,
} from '@influxdata/giraffe'
import {getFormatter} from 'src/utils/vis'
import {DEFAULT_LINE_COLORS} from 'src/constants/graphColorPalettes'
import {generateThresholdsListHexs} from 'src/constants/colorOperations'
import {DEFAULT_TIME_FORMAT} from 'src/constants/timeFormat'

interface Props extends VisualizationProps {
  properties: LinePlusSingleStatViewProperties
}

const SingleStatPlusLine: FunctionComponent<Props> = ({properties, result}) => {
  const {table} = result
  const timeFormat = properties.timeFormat || DEFAULT_TIME_FORMAT
  const xColumn = properties.xColumn || '_time'
  const yColumn =
    (table.columnKeys.includes(properties.yColumn) && properties.yColumn) ||
    '_value'

  const xFormatter = getFormatter('time', {
    prefix: properties.axes.x.prefix,
    suffix: properties.axes.x.suffix,
    base: properties.axes.x.base,
    timeZone: 'Local',
    timeFormat,
  })

  const yFormatter = getFormatter(table.getColumnType(yColumn), {
    prefix: properties.axes.y.prefix,
    suffix: properties.axes.y.suffix,
    base: properties.axes.y.base,
    timeZone: 'Local',
    timeFormat,
  })

  const groupKey = useMemo(() => [...result.fluxGroupKeyUnion], [result])
  const colorHexes = useMemo(() => {
    const _colors = properties.colors.filter(c => c.type === 'scale')
    if (_colors && _colors.length) {
      return _colors.map(color => color.hex)
    }
    return DEFAULT_LINE_COLORS.map(color => color.hex)
  }, [properties.colors])

  const latestValues = getLatestValues(result.table)
  const latestValue = latestValues[0]

  const {backgroundColor, textColor} = generateThresholdsListHexs({
    colors: properties.colors,
    lastValue: latestValue,
    cellType: 'single-stat',
  })

  const config: Config = {
    table,
    xAxisLabel: properties.axes.x.label,
    yAxisLabel: properties.axes.y.label,
    valueFormatters: {
      [xColumn]: xFormatter,
      [yColumn]: yFormatter,
    },
    layers: [
      {
        type: 'line',
        x: xColumn,
        y: yColumn,
        fill: groupKey,
        position: properties.position,
        colors: colorHexes,
        shadeBelow: !!properties.shadeBelow,
        shadeBelowOpacity: 0.08,
        hoverDimension: properties.hoverDimension,
      },
      {
        type: 'single stat',
        prefix: properties.prefix,
        suffix: properties.suffix,
        decimalPlaces: properties.decimalPlaces,
        textColor,
        backgroundColor: backgroundColor ? backgroundColor : '',
        svgTextStyle: {
          fontSize: '100',
          fontWeight: 'lighter',
          dominantBaseline: 'middle',
          textAnchor: 'middle',
          letterSpacing: '-0.05em',
        },
        svgStyle: SINGLE_STAT_SVG_NO_USER_SELECT,
      },
    ],
  }

  return <Plot config={config} />
}

export default SingleStatPlusLine
