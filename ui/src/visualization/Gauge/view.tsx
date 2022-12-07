import React, {FunctionComponent} from 'react'

import {GaugeViewProperties} from 'src/types/cells'
import {Config,  Plot } from '@influxdata/giraffe'
import {VisualizationProps} from 'src/visualization'
import { DEFAULT_GAUGE_COLORS } from 'src/constants/thresholds'
import { Color } from 'src/types/colors'

export const GAUGE_ARC_LENGTH_DEFAULT = 1.5 * Math.PI
export const GAUGE_VALUE_POSITION_X_OFFSET_DEFAULT = 0
export const GAUGE_VALUE_POSITION_Y_OFFSET_DEFAULT = 1.5

interface Props extends VisualizationProps {
  properties: GaugeViewProperties
}

const Gauge: FunctionComponent<Props> = ({result: {table}, properties}) => {
  const {prefix, tickPrefix, suffix, tickSuffix, decimalPlaces} =
    properties

  const config: Config = {
    table,
    layers: [
      {
        type: 'gauge',
        prefix,
        tickPrefix,
        suffix,
        tickSuffix,
        decimalPlaces,
        gaugeColors: DEFAULT_GAUGE_COLORS as Color[],
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
