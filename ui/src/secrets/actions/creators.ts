// Libraries
import {NormalizedSchema} from 'normalizr'

// Types
import {RemoteDataState, Sort} from '@influxdata/clockface'
import {SecretEntities} from 'src/types/schemas'
import {SortTypes} from 'src/types/sort'

export const SET_SECRET = 'SET_SECRET'
export const SET_SECRETS = 'SET_SECRETS'
export const REMOVE_SECRET = 'REMOVE_SECRET'
export const SET_SECRET_SEARCH_TERM = 'SET_SECRET_SEARCH_TERM'
export const SET_SECRET_SORT_OPTIONS = 'SET_SECRET_SORT_OPTIONS'

type SecretSchema<R extends string | string[]> = NormalizedSchema<
  SecretEntities,
  R
>

export const setSecret = (
  status: RemoteDataState,
  schema?: SecretSchema<string>
) =>
  ({
    type: SET_SECRET,
    id: schema.result,
    status,
    schema,
  } as const)

export const setSecrets = (
  status: RemoteDataState,
  schema?: SecretSchema<string[]>
) =>
  ({
    type: SET_SECRETS,
    status,
    schema,
  } as const)

export const removeSecret = (key: string) =>
  ({
    type: REMOVE_SECRET,
    id: key,
  } as const)

export const setSecretSearchTerm = (searchTerm: string) =>
  ({
    type: SET_SECRET_SEARCH_TERM,
    payload: searchTerm,
  } as const)

export const setSecretSortOptions = (
  key: string,
  direction: Sort,
  type: SortTypes
) =>
  ({
    type: SET_SECRET_SORT_OPTIONS,
    payload: {
      key,
      type,
      direction,
    },
  } as const)

export type Action =
  | ReturnType<typeof setSecret>
  | ReturnType<typeof setSecrets>
  | ReturnType<typeof removeSecret>
  | ReturnType<typeof setSecretSearchTerm>
  | ReturnType<typeof setSecretSortOptions>
