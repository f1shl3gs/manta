// Libraries
import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'
import moment from 'moment'

// Types
import {TimeRange} from '../types/TimeRanges'
import {AutoRefresh, AutoRefreshStatus} from 'types/AutoRefresh'

// Hooks
import {useTimeRange} from './useTimeRange'

const MAX_POINT = 1024
const MIN_STEP = 14

const calculateRange = (timeRange: TimeRange) => {
  switch (timeRange.type) {
    case 'selectable-duration':
      const end = moment().unix()
      const start = end - timeRange.seconds

      let step = (end - start) / MAX_POINT
      if (step < MIN_STEP) {
        step = MIN_STEP
      }

      return {
        start,
        end,
        step,
      }

    default:
      const now = moment().unix()

      return {
        start: now - 3600,
        end: now,
        step: 14,
      }
  }
}

const [AutoRefreshProvider, useAutoRefresh] = constate(
  () => {
    const {timeRange} = useTimeRange()

    const [autoRefresh, setAutoRefresh] = useState<AutoRefresh>({
      status: AutoRefreshStatus.Active,
      interval: 15,
    })
    const [state, setState] = useState(() => calculateRange(timeRange))
    // const [manualRefresh, setManualRefresh] = useState(0);
    const refresh = useCallback(() => {
      // setManualRefresh(prevState => prevState + 1);
      setState((prevState) => calculateRange(timeRange))
    }, [timeRange])

    useEffect(() => {
      setState((prevState) => calculateRange(timeRange))

      if (autoRefresh.status !== AutoRefreshStatus.Active) {
        return
      }

      const timer = setInterval(() => {
        if (document.hidden) {
          // no need to refresh
          return
        }

        setState((prev) => calculateRange(timeRange))
      }, autoRefresh.interval * 1000)

      return () => {
        clearInterval(timer)
      }
    }, [autoRefresh, timeRange])

    return {
      autoRefresh,
      setAutoRefresh,
      // manualRefresh,
      refresh,
      ...state,
    }
  },
  (value) => value
)

export {AutoRefreshProvider, useAutoRefresh}
