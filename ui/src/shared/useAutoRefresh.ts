import constate from 'constate'
import {useCallback, useState} from 'react'
import dayjs from 'dayjs'
import {useSearchParams} from 'react-router-dom'

import {PARAMS_INTERVAL, useTimeRange} from 'src/shared/useTimeRange'
import {AutoRefresh, AutoRefreshStatus} from 'src/types/autoRefresh'
import {parseDuration} from 'src/utils/duration'
import {TimeRange} from 'src/types/timeRanges'

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

  return {
    autoRefresh,
    setAutoRefresh,
    refresh,
    setRange: setState,
    ...state,
  }
})

export {AutoRefreshProvider, useAutoRefresh, calculateRange}
