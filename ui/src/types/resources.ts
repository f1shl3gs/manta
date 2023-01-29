import {RemoteDataState} from '@influxdata/clockface'
import {Dashboard} from 'src/types/dashboards'
import {Cell} from 'src/types/cells'
import {Organization} from 'src/types/organization'
import {Config} from 'src/types/config'
import {User} from 'src/types/user'
import {Scrape} from 'src/types/scrape'
import {Check, CheckSortParams} from 'src/types/checks'
import {Secret, SecretSortParams} from 'src/types/secrets'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'
import {SortOptions} from 'src/types/sort'

export interface Resource {
  id: string
  type: ResourceType
}

export enum ResourceType {
  Cells = 'cells',
  Checks = 'checks',
  Configs = 'configs',
  Dashboards = 'dashboards',
  Members = 'members',
  NotificationEndpoints = 'notificationEndpoints',
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

export interface ConfigsState extends NormalizedState<Config> {
  config: {
    status: RemoteDataState
    content: string
  }
}

export interface DashboardsState extends NormalizedState<Dashboard> {
  current: string

  searchTerm: string
  sortOptions: SortOptions
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

export interface NotificationEndpointState
  extends NormalizedState<NotificationEndpoint> {
  searchTerm: string
  sortOptions: SortOptions

  current: NotificationEndpoint
}

// ResourceState defines the types for normalized resources
export interface ResourceState {
  [ResourceType.Cells]: CellsState
  [ResourceType.Checks]: CheckState
  [ResourceType.Configs]: ConfigsState
  [ResourceType.Dashboards]: DashboardsState
  [ResourceType.Members]: MembersState
  [ResourceType.NotificationEndpoints]: NotificationEndpointState
  [ResourceType.Organizations]: OrgsState
  [ResourceType.Secrets]: SecretsState
  [ResourceType.Scrapes]: ScrapesState
}
