import {AppState as State} from 'src/shared/reducers/app'
import {AutoRefreshState} from 'src/shared/reducers/autoRefresh'
import {TimeRangeState} from 'src/shared/reducers/timeRange'
import {ResourceState} from 'src/types/resources'

// import {Notification} from 'src/types/Notification'

export interface AppState {
  app: State
  autoRefresh: AutoRefreshState

  // notifications: Notification[]

  resources: ResourceState

  timeRange: TimeRangeState
}

export type GetState = () => AppState
