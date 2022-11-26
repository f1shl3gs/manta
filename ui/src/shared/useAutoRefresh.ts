import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'
import dayjs from 'dayjs'
import {useSearchParams} from 'react-router-dom'

import {PARAMS_INTERVAL, useTimeRange} from 'src/shared/useTimeRange'
import {AutoRefresh, AutoRefreshStatus} from 'src/types/AutoRefresh'
import {parseDuration} from 'src/utils/duration'
import {TimeRange} from 'src/types/TimeRanges'

const MAX_POINT = 1024
const MIN_STEP = 14

const calculateRange = (timeRange: TimeRange) => {
  switch (timeRange.type) {
    case 'selectable-duration':
      const end = dayjs().unix()
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
      const now = dayjs().unix()

      return {
        start: now - 3600,
        end: now,
        step: 14,
      }
  }
}

const [AutoRefreshProvider, useAutoRefresh] = constate(() => {
  const {timeRange} = useTimeRange()
  const [params] = useSearchParams()
  const [autoRefresh, setAutoRefresh] = useState<AutoRefresh>(() => {
    const val = params.get(PARAMS_INTERVAL)
    if (val === null) {
      // TODO: set default and update params
      return {
        interval: 15,
        status: AutoRefreshStatus.Active,
      }
    }

    const duration = parseDuration(val)

    return {
      status:
        duration !== 0 ? AutoRefreshStatus.Active : AutoRefreshStatus.Paused,
      interval: duration / 1000,
    }
  })

  const [state, setState] = useState(() => calculateRange(timeRange))
  const refresh = useCallback(() => {
    setState(_prev => calculateRange(timeRange))
  }, [timeRange])

  useEffect(() => {
    setState(_ => calculateRange(timeRange))

    if (autoRefresh.status !== AutoRefreshStatus.Active) {
      return
    }

    const timer = setInterval(() => {
      if (document.hidden) {
        // tab is not focused, no need to refresh
        return
      }

      setState(_ => calculateRange(timeRange))
    }, autoRefresh.interval * 1000)

    return () => {
      clearInterval(timer)
    }
  }, [autoRefresh, timeRange])

  return {
    autoRefresh,
    setAutoRefresh,
    refresh,
    ...state,
  }
})

export {AutoRefreshProvider, useAutoRefresh, calculateRange}
