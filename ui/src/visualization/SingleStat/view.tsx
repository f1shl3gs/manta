import React, {FunctionComponent} from 'react'
import {VisualizationProps} from 'src/visualization'
import {SingleStatViewProperties} from 'src/types/cells'
import {Config, getLatestValues, Plot} from '@influxdata/giraffe'
import {generateThresholdsListHexs} from 'src/shared/constants/colorOperations'

interface Props extends VisualizationProps {
  properties: SingleStatViewProperties
}

const SingleStat: FunctionComponent<Props> = ({properties, result}) => {
  const {prefix, suffix, colors, decimalPlaces} = properties
  const {table} = result

  const latestValues = getLatestValues(table)
  const latestValue = latestValues[0]

  const {backgroundColor, textColor} = generateThresholdsListHexs({
    colors,
    lastValue: latestValue,
    cellType: 'single-stat',
  })

  const config: Config = {
    table,
    showAxes: false,
    layers: [
      {
        type: 'single stat',
        prefix,
        suffix,
        decimalPlaces,
        textColor,
        backgroundColor: backgroundColor ? backgroundColor : '',
        textOpacity: 100,
        svgTextStyle: {
          fontSize: '100',
          fontWeight: 'lighter',
          dominantBaseline: 'middle',
          textAnchor: 'middle',
          letterSpacing: '-0.05em',
        },
        svgTextAttributes: {
          'data-testid': 'single-stat--text',
        },
      },
    ],
  }

  return <Plot config={config} />
}

export default SingleStat
