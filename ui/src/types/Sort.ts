import {Sort} from '@influxdata/clockface'

export type SortKey = 'created' | 'updated' | 'name'

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
