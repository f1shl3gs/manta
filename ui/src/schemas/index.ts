// Libraries
import {schema} from "normalizr"
import {omit} from 'lodash'

// Types
import {ResourceType} from "src/types/resources"
import {Cell, Dashboard} from 'src/types/Dashboard'
import {RemoteDataState} from '@influxdata/clockface'

// Defines the schema for the 'organizations' resource
export const orgSchema = new schema.Entity(ResourceType.Organizations)
export const arrayOfOrgs = [orgSchema]

export const viewSchema = new schema.Entity(ResourceType.Views)
export const arrayOfViews = [viewSchema]

export const cellSchema = new schema.Entity(
  ResourceType.Cells,
  {},
  {
    processStrategy: (cell: Cell, parent: Dashboard) => {
      return {
        ...omit<Cell>(cell, 'properties'),
        dashboardID: cell.dashboardID ? cell.dashboardID : parent.id,
          status: RemoteDataState.Done
      }
    }
  }
)
export const arrayOfCells = [cellSchema]

// Defines the schema for the "dashboards" resource
export const dashboardSchema =new schema.Entity(
  ResourceType.Dashboards,
  {
    cells: arrayOfCells,
    views: arrayOfViews
  },
  {
    processStrategy: (dashboard: Dashboard) => addDashboardDefaults(dashboard)
  }
)
export const arrayOfDashboards = [dashboardSchema]

export const addDashboardDefaults = (dashboard: Dashboard): Dashboard => {
  return {
    ...dashboard,
    id: dashboard.id || '',
    name: dashboard.name || '',
    orgID: dashboard.orgID || '',
    status: RemoteDataState.Done
  }
}
