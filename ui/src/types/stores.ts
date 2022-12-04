import {AppState as State} from 'src/shared/reducers/app'
import {AutoRefreshState} from 'src/shared/reducers/autoRefresh'
import {TimeRangeState} from 'src/shared/reducers/timeRange'
import {ResourceState} from 'src/types/resources'
import {Notification} from 'src/types/notification'
import {ReduxRouterState} from '@lagunovsky/redux-react-router'
import {TimeMachineState} from 'src/timeMachine/reducers'

export interface AppState {
  router: ReduxRouterState

  app: State
  autoRefresh: AutoRefreshState

  notifications: Notification[]

  resources: ResourceState

  timeRange: TimeRangeState

  timeMachine: TimeMachineState
}

export type GetState = () => AppState