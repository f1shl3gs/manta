import {AppState as State} from 'src/shared/reducers/app'
import {AutoRefreshState} from 'src/shared/reducers/autoRefresh'
import {TimeRangeState} from 'src/shared/reducers/timeRange'
import {ResourceState} from 'src/types/resources'
import {Notification} from 'src/types/notification'
import {ReduxRouterState} from '@lagunovsky/redux-react-router'
import {TimeMachineState} from 'src/timeMachine/reducers'
import {MeState} from 'src/me/reducers'
import {SetupState} from 'src/setup/reducers'
import {CheckBuilderState} from 'src/checks/reducers/builder'

export interface AppState {
  router: ReduxRouterState

  app: State
  autoRefresh: AutoRefreshState
  checkBuilder: CheckBuilderState
  me: MeState
  notifications: Notification[]
  resources: ResourceState
  setup: SetupState
  timeRange: TimeRangeState
  timeMachine: TimeMachineState
}

export type GetState = () => AppState
