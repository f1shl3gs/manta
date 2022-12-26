import {produce} from 'immer'

import {TimeRange} from 'src/types/timeRanges'

import {pastHourTimeRange} from 'src/shared/constants/timeRange'

import {
  Action,
  ADD_QUERY,
  REMOVE_QUERY,
  RESET_TIMEMACHINE,
  SET_ACTIVE_QUERY,
  SET_ACTIVE_QUERY_TEXT,
  SET_QUERY_RESULTS,
  SET_VIEW_NAME,
  SET_VIEW_PROPERTIES,
  SET_VIEWING_VIS_OPTIONS,
  SET_VIEWTYPE,
} from 'src/timeMachine/actions'
import {ViewProperties} from 'src/types/cells'
import {RemoteDataState} from '@influxdata/clockface'
import {FromFluxResult, fromRows} from '@influxdata/giraffe'
import {createView} from 'src/visualization/helper'

export interface QueryResultsState {
  result: FromFluxResult
  state: RemoteDataState
}

export interface TimeMachineState {
  name: string
  activeQueryIndex: number | null
  viewProperties: ViewProperties
  viewingVisOptions: boolean
  contextID: string // dashboard, check or alert
  timeRange: TimeRange
  queryResult: QueryResultsState
}

const initialState = (): TimeMachineState => ({
  name: '',
  activeQueryIndex: 0,
  viewingVisOptions: false,
  contextID: '',
  timeRange: pastHourTimeRange,
  viewProperties: createView('xy'),
  queryResult: {
    state: RemoteDataState.NotStarted,
    result: {
      table: fromRows([]),
      fluxGroupKeyUnion: [],
      resultColumnNames: [],
    },
  },
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

      case SET_VIEW_NAME:
        draftState.name = action.name
        return

      case SET_VIEW_PROPERTIES:
        draftState.viewProperties = action.viewProperties
        return

      case RESET_TIMEMACHINE:
        return initialState()

      case SET_QUERY_RESULTS:
        draftState.queryResult = action.payload
        return

      case SET_VIEWTYPE:
        const newProperties = createView(action.viewType)
        draftState.viewProperties = {
          ...newProperties,
          queries: draftState.viewProperties.queries,
        }
        return

      default:
        return
    }
  })
