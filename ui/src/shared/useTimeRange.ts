// Libraries
import {useEffect, useState} from 'react'
import constate from 'constate'

// Types
import {TimeRange} from 'src/types/TimeRanges'
import {useSearchParams} from 'react-router-dom'
import {
  PARAMS_TIME_RANGE_LOW,
  PARAMS_TIME_RANGE_TYPE,
  pastHourTimeRange,
  SELECTABLE_TIME_RANGES,
} from 'src/constants/timeRange'

// Constants

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

  useEffect(() => {
    if (!params.get(PARAMS_TIME_RANGE_TYPE)) {
      return
    }

    setParams((prev: URLSearchParams) => {
      prev.set(PARAMS_TIME_RANGE_LOW, timeRange.lower)
      if (timeRange.upper) {
        prev.set('_upper', timeRange.upper)
      }

      return prev
    })
  }, [params, setParams, timeRange])

  return {
    timeRange,
    setTimeRange,
  }
})

export {TimeRangeProvider, useTimeRange}
