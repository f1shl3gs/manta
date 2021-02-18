// Libraries
import React from 'react'

// Components
import {AutoSizer} from 'react-virtualized'

import {useViewProperties} from '../../shared/useViewProperties'
import Gauge from './Gauge'
import {DashboardColor} from 'types/Colors'
import {GAUGE_THEME_DARK} from '../../constants/gauge'

const colors: DashboardColor[] = [
  {
    id: '0',
    type: 'min',
    hex: '#00C9FF',
    name: 'laser',
    value: 0,
  },
  {
    id: '608289e1-a646-47c8-8602-4b2dd14c5751',
    type: 'threshold',
    hex: '#FFB94A',
    name: 'pineapple',
    value: 20,
  },
  {
    id: 'a2ae0345-1189-43de-8cf1-8aaa441bebd9',
    type: 'threshold',
    hex: '#DC4E58',
    name: 'fire',
    value: 85,
  },
  {
    id: '1',
    type: 'max',
    hex: '#9394FF',
    name: 'comet',
    value: 100,
  },
]

interface Props {
  value: number
}

const GaugeChart: React.FC<Props> = props => {
  const {value} = props

  return (
    <AutoSizer>
      {({width, height}) => (
        <div className={'gauge'}>
          <Gauge
            width={width}
            height={height}
            colors={colors}
            prefix={''}
            suffix={''}
            tickPrefix={''}
            tickSuffix={''}
            gaugePosition={value}
            decimalPlaces={{
              isEnforced: true,
              digits: 2,
            }}
            theme={GAUGE_THEME_DARK}
          />
        </div>
      )}
    </AutoSizer>
  )
}

export default GaugeChart
