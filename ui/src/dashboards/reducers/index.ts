import {produce} from 'immer'

import {ResourceState, ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {DEFAULT_DASHBOARD_SORT_OPTIONS} from 'src/constants/dashboard'

import {
  Action as CellAction,
  SET_CELL,
  SET_CELLS,
} from 'src/cells/actions/creators'

import {
  Action,
  EDIT_DASHBOARD,
  REMOVE_DASHBOARD,
  SET_DASHBOARD,
  SET_DASHBOARD_SEARCH_TERM,
  SET_DASHBOARD_SORT,
  SET_DASHBOARDS,
} from 'src/dashboards/actions/creators'
import {
  editResource,
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'
import {Dashboard} from 'src/types/dashboards'

type DashboardsState = ResourceState['dashboards']

const initialState = (): DashboardsState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  sortOptions: DEFAULT_DASHBOARD_SORT_OPTIONS,
  searchTerm: '',
  current: '',
})

export const dashboardsReducer = (
  state: DashboardsState = initialState(),
  action: Action | CellAction
): DashboardsState => {
  return produce(state, draftState => {
    switch (action.type) {
      case SET_DASHBOARDS:
        setResource<Dashboard>(draftState, action, ResourceType.Dashboards)
        return

      case REMOVE_DASHBOARD:
        removeResource<Dashboard>(draftState, action)
        return

      case EDIT_DASHBOARD:
        editResource<Dashboard>(draftState, action, ResourceType.Dashboards)
        return

      case SET_DASHBOARD:
        draftState.current = action.id
        setResourceAtID<Dashboard>(draftState, action, ResourceType.Dashboards)
        return

      case SET_DASHBOARD_SORT:
        draftState.sortOptions = action.sortOptions
        return

      case SET_DASHBOARD_SEARCH_TERM:
        draftState.searchTerm = action.searchTerm
        return

      case SET_CELL:
        const {schema} = action

        const cellID = schema.result
        const cell = schema.entities.cells[cellID]
        const dashboards = draftState.byID[cell.dashboardID]

        if (dashboards?.cells.includes(cellID)) {
          return
        }

        if (draftState.byID[cell.dashboardID]) {
          draftState.byID[cell.dashboardID].cells.push(cellID)
        }

        return

      case SET_CELLS: {
        const {dashboardID, schema} = action
        const cellIDs = schema && schema.result

        if (!cellIDs) {
          return
        }

        draftState.byID[dashboardID].cells = cellIDs
        return
      }

      default:
        return
    }
  })
}
