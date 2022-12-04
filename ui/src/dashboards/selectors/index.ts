// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {Dashboard} from 'src/types/dashboards'

export const getDashboard = (state: AppState): Dashboard => {
  return get(state, 'resources.dashboards.dashboard', null)
}
