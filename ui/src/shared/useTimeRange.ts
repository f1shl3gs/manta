// Libraries
import {useState} from 'react'
import constate from 'constate'

// Types
import {TimeRange} from 'src/types/timeRanges'
import {SelectableDurationTimeRange} from 'src/types/timeRanges'
import {useSearchParams} from 'react-router-dom'

// Constants
export const PARAMS_INTERVAL = 'interval'
export const PARAMS_TIME_RANGE_LOW = 'low'
export const PARAMS_TIME_RANGE_TYPE = '_type'
export const PARAMS_SHOW_VARIABLES_CONTROLS = 'showVariablesControls'

export const TIME_RANGE_FORMAT = 'YYYY-MM-DD HH:mm'

export const CUSTOM_TIME_RANGE_LABEL = 'Custom Time Range'

export const pastFifteenMinTimeRange: SelectableDurationTimeRange = {
  seconds: 900,
  lower: 'now() - 15m',
  upper: null,
  label: 'Past 15m',
  duration: '15m',
  type: 'selectable-duration',
  windowPeriod: 10000, // 10s
}

export const pastHourTimeRange: SelectableDurationTimeRange = {
  seconds: 3600,
  lower: 'now() - 1h',
  upper: null,
  label: 'Past 1h',
  duration: '1h',
  type: 'selectable-duration',
  windowPeriod: 10000, // 10s
}

export const pastThirtyDaysTimeRange: SelectableDurationTimeRange = {
  seconds: 2592000,
  lower: 'now() - 30d',
  upper: null,
  label: 'Past 30d',
  duration: '30d',
  type: 'selectable-duration',
  windowPeriod: 3600000, // 1h
}

export const pastThirtyMinutesTimeRange: SelectableDurationTimeRange = {
  seconds: 1800,
  lower: 'now() - 30m',
  upper: null,
  label: 'Past 30m',
  duration: '30m',
  type: 'selectable-duration',
  windowPeriod: 10,
}

export const SELECTABLE_TIME_RANGES: SelectableDurationTimeRange[] = [
  {
    seconds: 300,
    lower: 'now() - 5m',
    upper: null,
    label: 'Past 5m',
    duration: '5m',
    type: 'selectable-duration',
    windowPeriod: 10000, // 10s
  },
  pastFifteenMinTimeRange,
  pastThirtyMinutesTimeRange,
  pastHourTimeRange,
  {
    seconds: 21600,
    lower: 'now() - 6h',
    upper: null,
    label: 'Past 6h',
    duration: '6h',
    type: 'selectable-duration',
    windowPeriod: 60000, // 1m
  },
  {
    seconds: 43200,
    lower: 'now() - 12h',
    upper: null,
    label: 'Past 12h',
    duration: '12h',
    type: 'selectable-duration',
    windowPeriod: 120000, // 2m
  },
  {
    seconds: 86400,
    lower: 'now() - 24h',
    upper: null,
    label: 'Past 24h',
    duration: '24h',
    type: 'selectable-duration',
    windowPeriod: 240000, // 4m
  },
  {
    seconds: 172800,
    lower: 'now() - 2d',
    upper: null,
    label: 'Past 2d',
    duration: '2d',
    type: 'selectable-duration',
    windowPeriod: 600000, // 10m
  },
  {
    seconds: 604800,
    lower: 'now() - 7d',
    upper: null,
    label: 'Past 7d',
    duration: '7d',
    type: 'selectable-duration',
    windowPeriod: 1800000, // 30 min
  },
  pastThirtyDaysTimeRange,
]

const [TimeRangeProvider, useTimeRange] = constate(() => {
  const [params, setParams] = useSearchParams()
  const [timeRange, setTimeRange] = useState<TimeRange>(() => {
    switch (params.get(PARAMS_TIME_RANGE_TYPE)) {
      case 'selectable-duration': {
        const lower = params.get(PARAMS_TIME_RANGE_TYPE)

        return (
          SELECTABLE_TIME_RANGES.find(tr => lower === tr.lower) ||
          pastHourTimeRange
        )
      }
      case 'duration':
        return pastHourTimeRange
      case 'custom':
        return pastHourTimeRange
      default:
        return pastHourTimeRange
    }
  })

  const setTimeRangeWrapper = (tr: TimeRange) => {
    setTimeRange(tr)

    setParams(params => {
      if (!params.get(PARAMS_TIME_RANGE_TYPE)) {
        return params
      }

      params.set(PARAMS_TIME_RANGE_LOW, tr.lower)
      if (tr.upper) {
        params.set('_upper', tr.upper)
      }

      return params
    }, {replace: true})
  }

  return {
    timeRange,
    setTimeRange: setTimeRangeWrapper,
  }
})

export {TimeRangeProvider, useTimeRange}
