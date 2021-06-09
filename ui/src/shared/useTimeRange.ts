// Libraries
import {useEffect, useState} from 'react'
import constate from 'constate'

// Types
import {TimeRange} from '../types/TimeRanges'

// Hooks
import useSearchParams from './useSearchParams'

// Constants
import {pastHourTimeRange, SELECTABLE_TIME_RANGES} from '../constants/timeRange'
import {
  PARAMS_TIME_RANGE_LOW,
  PARAMS_TIME_RANGE_TYPE,
} from '../constants/params'

const [TimeRangeProvider, useTimeRange] = constate(
  () => {
    const {params, setParams} = useSearchParams()
    const [timeRange, setTimeRange] = useState<TimeRange>(() => {
      switch (params.get(PARAMS_TIME_RANGE_TYPE)) {
        case 'selectable-duration':
          const lower = params.get('_type')
          return (
            SELECTABLE_TIME_RANGES.find(tr => lower === tr.lower) ||
            pastHourTimeRange
          )
        case 'duration':
          return pastHourTimeRange
        case 'custom':
          return pastHourTimeRange
        default:
          return pastHourTimeRange
      }
    })

    useEffect(() => {
      if (!params.get('_type')) {
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
  },
  value => value
)

export {TimeRangeProvider, useTimeRange}
