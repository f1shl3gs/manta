import {ResourceState, ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {DEFAULT_DASHBOARD_SORT_OPTIONS} from 'src/constants/dashboard'
import {produce} from 'immer'

import {Action} from 'src/dashboards/actions/creators'
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'
import {Dashboard} from 'src/types/Dashboard'

type DashboardsState = ResourceState['dashboards']

const initialState = (): DashboardsState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  searchTerm: '',
  sortOptions: DEFAULT_DASHBOARD_SORT_OPTIONS,
})

export const dashboardsReducer = (
  state: DashboardsState = initialState(),
  action: Action
): DashboardsState => {
  return produce(state, draftState => {
    switch (action.type) {
      case 'SET_DASHBOARDS':
        setResource<Dashboard>(draftState, action, ResourceType.Dashboards)
        return

      case 'REMOVE_DASHBOARD':
        removeResource<Dashboard>(draftState, action)
        return

      case 'SET_DASHBOARD':
        setResourceAtID<Dashboard>(draftState, action, ResourceType.Dashboards)
        return

      case 'SET_DASHBOARD_SORT':
        draftState.sortOptions = action.sortOptions
        return

      case 'SET_DASHBOARD_SEARCH_TERM':
        draftState.searchTerm = action.searchTerm
        return

      default:
        return
    }
  })
}
