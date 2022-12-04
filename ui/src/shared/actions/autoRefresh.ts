import dayjs from 'dayjs'

import {AutoRefreshStatus} from 'src/types/autoRefresh'
import {TimeRange} from 'src/types/timeRanges'

export const SET_AUTOREFRESH_INTERVAL = 'SET_AUTOREFRESH_INTERVAL'
export const POLL = 'POLL'

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

export const setAutoRefreshInterval = (interval: number) =>
  ({
    type: SET_AUTOREFRESH_INTERVAL,
    payload: {
      interval,
      status:
        interval !== 0 ? AutoRefreshStatus.Active : AutoRefreshStatus.Paused,
    },
  } as const)

export const poll = (timeRange: TimeRange) => {
  const range = calculateRange(timeRange)

  return {
    type: POLL,
    payload: range,
  } as const
}

export type Action =
  | ReturnType<typeof poll>
  | ReturnType<typeof setAutoRefreshInterval>
