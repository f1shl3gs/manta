import {NormalizedSchema} from 'normalizr'

import {RemoteDataState} from '@influxdata/clockface'

import {ConfigurationEntities} from 'src/types/schemas'

// Action Types
export const SET_CONFIGS = 'SET_CONFIGS'
export const SET_CONFIG = 'SET_CONFIG'
export const REMOVE_CONFIG = 'REMOVE_CONFIG'

export type Action =
  | ReturnType<typeof setConfigs>
  | ReturnType<typeof removeConfig>
  | ReturnType<typeof setConfig>

export const setConfigs = (
  status: RemoteDataState,
  schema?: NormalizedSchema<ConfigurationEntities, string[]>
) =>
  ({
    type: SET_CONFIGS,
    status,
    schema,
  } as const)

export const removeConfig = (id: string) =>
  ({
    type: REMOVE_CONFIG,
    id,
  } as const)

export const setConfig = (
  id: string,
  status: RemoteDataState,
  schema?: NormalizedSchema<ConfigurationEntities, string>
) =>
  ({
    type: SET_CONFIG,
    id,
    schema,
    status,
  } as const)
