import {RemoteDataState} from '@influxdata/clockface'
import {Dashboard, DashboardSortParams} from 'src/types/Dashboard'
import {Organization} from 'src/types/Organization'

export enum ResourceType {
  Organizations = 'organizations',

  Dashboards = 'dashboards',
  Users = 'users',
  Configurations = 'configurations',

  Scrapes = 'scrapes',

  Cells = 'cells',
  Views = 'views'
}

export interface NormalizedState<R> {
  byID: {
    [uuid: string]: R
  }
  allIDs: string[]
  status: RemoteDataState
}

export interface DashboardsState extends NormalizedState<Dashboard> {
  searchTerm: string
  sortOptions: DashboardSortParams
}

export interface OrgsState extends NormalizedState<Organization> {
  org: Organization
}

// ResourceState defines the types for normalized resources
export interface ResourceState {
  [ResourceType.Dashboards]: DashboardsState
  [ResourceType.Organizations]: OrgsState
}
