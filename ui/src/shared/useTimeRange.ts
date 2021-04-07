// Libraries
import {useEffect, useState} from 'react'
import constate from 'constate'

// Types
import {
  SelectableDurationTimeRange,
  SelectableTimeRangeLower,
  TimeRange,
} from '../types/TimeRanges'

// Hooks
import useSearchParams from './useSearchParams'

// Constants
import {pastHourTimeRange} from '../constants/timeRange'

const [TimeRangeProvider, useTimeRange] = constate(
  () => {
    const {params, setParams} = useSearchParams()
    const [timeRange, setTimeRange] = useState<TimeRange>(() => {
      switch (params.get('_type')) {
        case 'selectable-duration':
          const tr = {} as SelectableDurationTimeRange

          const lower = params.get('_lower')
          if (lower) {
            tr.lower = lower as SelectableTimeRangeLower
          }

          return tr
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
        prev.set('_lower', timeRange.lower)
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
