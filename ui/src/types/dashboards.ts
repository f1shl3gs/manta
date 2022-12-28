import {RemoteDataState, Sort} from '@influxdata/clockface'
import {SortTypes} from 'src/types/sort'
import {ViewProperties} from 'src/types/cells'

export interface GenView {
  readonly id?: string
  name: string
  properties: ViewProperties
}

export interface View<T extends ViewProperties = ViewProperties>
  extends GenView {
  properties: T
  cellID?: string
  dashboardID?: string
}

export type NewView<T extends ViewProperties = ViewProperties> = Omit<
  View<T>,
  'id'
>

export interface Dashboard {
  readonly id: string
  readonly created: string
  readonly updated: string
  name: string
  desc: string
  orgID: string
  cells: string[]

  // TODO: remove this
  status: RemoteDataState
}

export interface DashboardQuery {
  name?: string
  text: string
  hidden: boolean
}

export type Dashboards = Dashboard[]

export interface DashboardSortParams {
  direction: Sort
  type: SortTypes
  key: string
}
