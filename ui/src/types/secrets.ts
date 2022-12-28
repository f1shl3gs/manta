import {RemoteDataState, Sort} from '@influxdata/clockface'
import {SortTypes} from 'src/types/sort'

export interface Secret {
  readonly key: string
  readonly updated: string

  status: RemoteDataState
}

export type SecretSortKey = keyof Secret

export interface SecretSortParams {
  direction: Sort
  type: SortTypes
  key: SecretSortKey
}
