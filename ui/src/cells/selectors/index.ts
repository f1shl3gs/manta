import {createSelector} from 'reselect'

import {Cell} from 'src/types/cells'
import {AppState} from 'src/types/stores'

const getResources = (state: AppState) => state.resources
const getDashboardId = (_, dashbaordID) => dashbaordID

export const getCells = createSelector(
  getResources,
  getDashboardId,
  (resources, dashboardID): Cell[] => {
    const dashbaord = resources.dashboards.byID[dashboardID]
    if (!dashbaord || !dashbaord.cells) {
      return []
    }

    return dashbaord.cells
      .filter(id => Boolean(resources.cells.byID[id]))
      .map(id => resources.cells.byID[id])
  }
)
