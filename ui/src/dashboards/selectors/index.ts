// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'

export const getDashboardID = (state: AppState): string => {
  return get(state, 'resources.dashboards.current', '')
}
