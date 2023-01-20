// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {getByID} from '../../resources/selectors'
import {Dashboard} from '../../types/dashboards'

export const getDashboardID = (state: AppState): string => {
  return get(state, 'resources.dashboards.current', '')
}

export const getDashboard = (id: string) => (state: AppState) => {
  return getByID<Dashboard>(state, ResourceType.Dashboards, id)
}

export const getDashboardWithCell = (id: string) => (state: AppState) => {
  const dashboard = getByID<Dashboard>(state, ResourceType.Dashboards, id)
  const cells = dashboard.cells.map(id => {
    return getByID(state, ResourceType.Cells, id)
  })

  return {
    ...dashboard,
    cells,
  }
}
