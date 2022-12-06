import {AppState} from 'src/types/stores'
import {DashboardQuery} from 'src/types/dashboards'
import {ViewProperties} from 'src/types/cells'

export const getTimeMachine = (state: AppState) => {
  return state.timeMachine
}

export const getQueries = (state: AppState): DashboardQuery[] => {
  return state.timeMachine.viewProperties.queries
}

export const getViewProperties = (state: AppState): ViewProperties =>
  state.timeMachine.viewProperties

export const getActiveQuery = (state: AppState) => {
  const activeIndex = state.timeMachine.activeQueryIndex
  const queries = getQueries(state)

  return queries[activeIndex]
}
