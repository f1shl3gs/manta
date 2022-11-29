import {Cell, Dashboard} from 'src/types/Dashboard'
import {Organization} from 'src/types/Organization'

// DashboardEntities defines the result of normalizr's normalization of the
// "dashboards" resource
export interface DashboardEntities {
  dashboards: {
    [uuid: string]: Dashboard
  }
  cells: {
    [uuid: string]: Cell
  }
}

// OrgEntities defines the result of normalizr's normalization
// of the "organizations" resource
export interface OrgEntities {
  orgs: {
    [uuid: string]: Organization
  }
}
