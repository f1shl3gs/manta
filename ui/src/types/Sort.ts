import {Sort} from '@influxdata/clockface'

import {Dashboard} from 'src/types/Dashboard'

export type DashboardSortKey = keyof Dashboard

export type SortKey = DashboardSortKey

export enum SortTypes {
  String = 'string',
  Date = 'date',
  Float = 'float',
}

export interface SortOption {
  key: SortKey
  type: SortTypes
  direction: Sort
}
