// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {Dashboard} from 'src/types/dashboards'

// Selectors
import {getByID} from 'src/resources/selectors'

// Utils
import {get} from 'src/shared/utils/get'

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
