import {RemoteDataState} from '@influxdata/clockface'
import {Dashboard, DashboardSortParams} from 'src/types/dashboards'
import {Cell} from 'src/types/cells'
import {Organization} from 'src/types/organization'
import {Configuration} from 'src/types/configuration'
import {User} from 'src/types/user'

export enum ResourceType {
  Cells = 'cells',
  Dashboards = 'dashboards',
  Configurations = 'configurations',
  Organizations = 'organizations',
  Scrapes = 'scrapes',
  Users = 'users',
  Views = 'views',
}

export interface NormalizedState<R> {
  byID: {
    [uuid: string]: R
  }
  allIDs: string[]
  status: RemoteDataState
}

export interface ConfigurationsState extends NormalizedState<Configuration> {
  config: {
    status: RemoteDataState
    content: string
  }
}

export interface DashboardsState extends NormalizedState<Dashboard> {
  searchTerm: string
  sortOptions: DashboardSortParams
}

export interface OrgsState extends NormalizedState<Organization> {
  org: Organization
}

export interface MembersState extends NormalizedState<User> {
  me: User
}

// Cells 'allIDs' are Dashboard.cells
type CellsState = Omit<NormalizedState<Cell>, 'allIDs'>

// ResourceState defines the types for normalized resources
export interface ResourceState {
  [ResourceType.Cells]: CellsState
  [ResourceType.Configurations]: ConfigurationsState
  [ResourceType.Dashboards]: DashboardsState
  [ResourceType.Organizations]: OrgsState
  [ResourceType.Members]: MembersState
}
