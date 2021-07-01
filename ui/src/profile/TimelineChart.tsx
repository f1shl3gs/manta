// Libraries
import React from 'react'
import {Config, newTable, Plot} from '@influxdata/giraffe'

// TestData
import testData from './components/testData'
import {getFormatter} from '../utils/vis'
import {DEFAULT_TIME_FORMAT} from '../constants/timeFormat'

const TimelineChart: React.FC = () => {
  const {
    timeline: {startTime, samples, durationDelta},
  } = testData

  const table = newTable(samples.length)
    .addColumn(
      'time',
      'dateTime:RFC3339',
      'time',
      samples.map((v, index) => startTime * 1000 + index * durationDelta)
    )
    .addColumn('samples', 'double', 'number', samples)

  const xFormatter = getFormatter('time', {
    base: '10',
    timeZone: 'Local',
    timeFormat: DEFAULT_TIME_FORMAT,
  })

  const config: Config = {
    table,
    valueFormatters: {
      // @ts-ignore
      _time: xFormatter,
    },
    layers: [
      {
        type: 'line',
        x: 'time',
        y: 'samples',
      },
    ],
  }

  return (
    <div style={{height: '200px'}}>
      <Plot config={config} />
    </div>
  )
}

export default TimelineChart
