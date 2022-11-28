export type Action = SetAutoRefresh

interface SetAutoRefresh {
  type: 'SetAutoRefreshInterval'
  payload: {
    second: number
  }
}

export const setAutoRefreshInterval = (second: number): SetAutoRefresh => ({
  type: 'SetAutoRefreshInterval',
  payload: {
    second,
  },
})
