// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from '../../types/resources'

export const getDashboardID = (state: AppState): string => {
  return get(state, 'resources.dashboards.current', '')
}

export const getDashboard = (id: string) => (state: AppState) => {
  return state.resources[ResourceType.Dashboards]['byID'][id]
}
