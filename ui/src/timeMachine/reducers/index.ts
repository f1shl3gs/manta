import {produce} from 'immer'

import {TimeRange} from 'src/types/timeRanges'

import {pastHourTimeRange} from 'src/constants/timeRange'

import {
  Action,
  ADD_QUERY,
  REMOVE_QUERY,
  SET_ACTIVE_QUERY,
  SET_ACTIVE_QUERY_TEXT,
  SET_VIEWING_VIS_OPTIONS,
} from 'src/timeMachine/actions'
import {DEFAULT_VIEWPROPERTIES} from 'src/constants/dashboard'
import {ViewProperties} from 'src/types/cells'

export interface TimeMachineState {
  name: string
  activeQueryIndex: number | null
  viewProperties: ViewProperties
  viewingVisOptions: boolean
  contextID: string // dashboard, check or alert
  timeRange: TimeRange
}

const initialState = () => ({
  name: '',
  activeQueryIndex: 0,
  queries: [
    {
      name: 'query 1',
      text: '',
      hidden: false,
    },
  ],
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

      case SET_ACTIVE_QUERY_TEXT:
        draftState.viewProperties.queries[draftState.activeQueryIndex].text =
          action.text
        return

      case ADD_QUERY:
        draftState.viewProperties.queries.push(action.query)
        return

      case SET_VIEWING_VIS_OPTIONS:
        draftState.viewingVisOptions = action.viewingVisOptions
        return

      case REMOVE_QUERY:
        draftState.viewProperties.queries.filter(
          (_q, index) => index === action.index
        )
        return

      default:
        return
    }
  })
