import {produce} from 'immer'

import {DashboardQuery, ViewProperties} from 'src/types/dashboards'
import {TimeRange} from 'src/types/timeRanges'

import {pastHourTimeRange} from 'src/constants/timeRange'

import {
  Action,
  ADD_QUERY,
  REMOVE_QUERY,
  SET_ACTIVE_QUERY,
  SET_VIEWING_VIS_OPTIONS,
} from 'src/timeMachine/actions'
import {DEFAULT_VIEWPROPERTIES} from 'src/constants/dashboard'

export interface TimeMachineState {
  activeQueryIndex: number | null
  queries: DashboardQuery[]
  viewProperties: ViewProperties
  viewingVisOptions: boolean
  contextID: string // dashboard, check or alert
  timeRange: TimeRange
}

const initialState = () => ({
  activeQueryIndex: null,
  queries: new Array<DashboardQuery>(),
  viewingVisOptions: false,
  contextID: '',
  timeRange: pastHourTimeRange,
  viewProperties: DEFAULT_VIEWPROPERTIES,
})

export const timeMachineReducer = (
  state: TimeMachineState = initialState(),
  action: Action
): TimeMachineState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_ACTIVE_QUERY:
        draftState.activeQueryIndex = action.activeQueryIndex
        return
      case ADD_QUERY:
        draftState.queries.push(action.query)
        return
      case SET_VIEWING_VIS_OPTIONS:
        draftState.viewingVisOptions = action.viewingVisOptions
        return
      case REMOVE_QUERY:
        draftState.queries.filter((_q, index) => index === action.index)
        return
      default:
        return
    }
  })
