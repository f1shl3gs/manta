import {RemoteDataState} from '@influxdata/clockface'
import {Dashboard, DashboardSortParams} from 'src/types/dashboards'
import {Cell} from 'src/types/cells'
import {Organization} from 'src/types/organization'
import {Configuration} from 'src/types/configuration'
import {User} from 'src/types/user'
import {Scrape} from 'src/types/scrape'
import {Check, CheckSortParams} from 'src/types/checks'
import {Secret, SecretSortParams} from 'src/types/secrets'

export interface Resource {
  type: ResourceType
  id: string
}

export enum ResourceType {
  Cells = 'cells',
  Checks = 'checks',
  Configurations = 'configurations',
  Dashboards = 'dashboards',
  Members = 'members',
  Organizations = 'organizations',
  Secrets = 'secrets',
  Scrapes = 'scrapes',
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
  current: string

  searchTerm: string
  sortOptions: DashboardSortParams
}

export interface OrgsState extends NormalizedState<Organization> {
  org: Organization
}

export type MembersState = NormalizedState<User>

export type ScrapesState = NormalizedState<Scrape>

export interface SecretsState extends NormalizedState<Secret> {
  searchTerm: string
  sortOptions: SecretSortParams
}

// Cells 'allIDs' are Dashboard.cells
export type CellsState = Omit<NormalizedState<Cell>, 'allIDs'>

export interface CheckState extends NormalizedState<Check> {
  searchTerm: string
  sortOptions: CheckSortParams
}

// ResourceState defines the types for normalized resources
export interface ResourceState {
  [ResourceType.Cells]: CellsState
  [ResourceType.Checks]: CheckState
  [ResourceType.Configurations]: ConfigurationsState
  [ResourceType.Dashboards]: DashboardsState
  [ResourceType.Organizations]: OrgsState
  [ResourceType.Members]: MembersState
  [ResourceType.Secrets]: SecretsState
  [ResourceType.Scrapes]: ScrapesState
}
