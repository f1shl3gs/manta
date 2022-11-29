import {AppState as State} from 'src/shared/reducers/app'
import {AutoRefreshState} from 'src/shared/reducers/autoRefresh'
import {TimeRangeState} from 'src/shared/reducers/timeRange'
// import {Notification} from 'src/types/Notification'

export interface AppState {
  app: State
  autoRefresh: AutoRefreshState

  // notifications: Notification[]

  timeRange: TimeRangeState
}

export type GetState = () => AppState
