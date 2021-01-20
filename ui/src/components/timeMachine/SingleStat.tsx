// Libraries
import React from 'react'

// Constants
import {generateThresholdsListHexs} from 'constants/colorOperations'

// Types
import {SingleStatViewProperties} from '../../types/Dashboard'

import {formatStatValue} from './Gauge'

interface Props {
  stat: number
  theme: 'light' | 'dark'
  properties: SingleStatViewProperties
}

const SingleStat: React.FC<Props> = (props) => {
  const {stat, theme, properties} = props
  const {colors, prefix, suffix, decimalPlaces} = properties

  const {bgColor: backgroundColor, textColor} = generateThresholdsListHexs({
    colors,
    lastValue: stat,
    cellType: 'single-stat',
  })

  const formattedValue = formatStatValue(stat, {decimalPlaces, prefix, suffix})
  return (
    <div className={'single-stat'} style={{backgroundColor}}>
      <svg
        width={'100%'}
        height={'100%'}
        viewBox={`0 0 ${formattedValue.length * 55} 100`}
      >
        <text
          className={'single-stat--text'}
          fontSize={'100'}
          y={'59%'}
          x={'50%'}
          dominantBaseline={'middle'}
          textAnchor={'middle'}
          style={{fill: textColor}}
        >
          {formattedValue}
        </text>
      </svg>
    </div>
  )
}

export default SingleStat
