import React, {FunctionComponent} from 'react'

import {GaugeViewProperties} from 'src/types/dashboards'
import {Config, Plot} from '@influxdata/giraffe'
import {VisualizationProps} from 'src/visualization'

export const GAUGE_ARC_LENGTH_DEFAULT = 1.5 * Math.PI
export const GAUGE_VALUE_POSITION_X_OFFSET_DEFAULT = 0
export const GAUGE_VALUE_POSITION_Y_OFFSET_DEFAULT = 1.5

interface Props extends VisualizationProps {
  properties: GaugeViewProperties
}

const Gauge: FunctionComponent<Props> = props => {
  const {colors, prefix, suffix, decimalPlaces} = props.properties
  const {table} = props.result

  const config: Config = {
    table,
    layers: [
      {
        type: 'gauge',
        prefix,
        suffix,
        decimalPlaces,
        gaugeColors: colors,
        gaugeSize: GAUGE_ARC_LENGTH_DEFAULT,
        gaugeTheme: {
          valuePositionXOffset: GAUGE_VALUE_POSITION_X_OFFSET_DEFAULT,
          valuePositionYOffset: GAUGE_VALUE_POSITION_Y_OFFSET_DEFAULT,
        },
      },
    ],
  }

  return <Plot config={config} />
}

export default Gauge
