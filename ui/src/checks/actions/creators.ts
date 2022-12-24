// Libraries
import {RemoteDataState} from '@influxdata/clockface'
import {NormalizedSchema} from 'normalizr'

// Types
import {CheckEntities} from 'src/types/schemas'

// Actions
export const SET_CHECKS = 'SET_CHECKS'
export const SET_CHECK = 'SET_CHECK'
export const REMOVE_CHECK = 'REMOVE_CHECK'
export const SET_CHECK_SEARCH_TERM = 'SET_CHECK_SEARCH_TERM'

type CheckSchema<R extends string | string[]> = NormalizedSchema<
  CheckEntities,
  R
>

export const setChecks = (
  status: RemoteDataState,
  schema?: CheckSchema<string[]>
) =>
  ({
    type: SET_CHECKS,
    status,
    schema,
  } as const)

export const setCheck = (
  id: string,
  status: RemoteDataState,
  schema?: CheckSchema<string>
) =>
  ({
    type: SET_CHECK,
    id,
    status,
    schema,
  } as const)

export const removeCheck = (id: string) =>
  ({
    type: REMOVE_CHECK,
    id,
  } as const)

export const setCheckSearchTerm = (searchTerm: string) =>
  ({
    type: SET_CHECK_SEARCH_TERM,
    searchTerm,
  } as const)

export type Action =
  | ReturnType<typeof setChecks>
  | ReturnType<typeof setCheck>
  | ReturnType<typeof removeCheck>
  | ReturnType<typeof setCheckSearchTerm>
