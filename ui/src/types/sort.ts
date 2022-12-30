import {Sort} from '@influxdata/clockface'

export enum SortTypes {
  String = 'string',
  Date = 'date',
  Float = 'float',
}

export interface SortOptions {
  key: string
  direction: Sort
  type: SortTypes
}
