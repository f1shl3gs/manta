import {NormalizedSchema} from 'normalizr'
import {CellEntities} from 'src/types/cells'
import {RemoteDataState} from '@influxdata/clockface'

// Action types
export const EDIT_CELL = 'EDIT_CELL'
export const SET_CELLS = 'SET_CELLS'
export const SET_CELL = 'SET_CELL'
export const REMOVE_CELL = 'REMOVE_CELL'

// R is the type of the value of the 'result' key in normalizr's normalization
export type CellSchema<R extends string | string[]> = NormalizedSchema<
  CellEntities,
  R
>

export const removeCell = (dashboardID: string, id: string) =>
  ({
    type: REMOVE_CELL,
    dashboardID,
    id,
  } as const)

export const setCells = (
  dashboardID: string,
  status: RemoteDataState,
  schema?: CellSchema<string[]>
) =>
  ({
    type: SET_CELLS,
    dashboardID,
    status,
    schema,
  } as const)

export const setCell = (
  id: string,
  status: RemoteDataState,
  schema?: CellSchema<string>
) =>
  ({
    type: SET_CELL,
    id,
    status,
    schema,
  } as const)

export const editCell = (schema: CellSchema<string>) =>
  ({
    type: EDIT_CELL,
    schema,
  } as const)

export type Action =
  | ReturnType<typeof removeCell>
  | ReturnType<typeof setCells>
  | ReturnType<typeof setCell>
  | ReturnType<typeof editCell>
