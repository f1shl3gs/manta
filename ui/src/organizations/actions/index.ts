import {RemoteDataState} from '@influxdata/clockface'
import {Organization} from 'src/types/organization'
import {NormalizedSchema} from 'normalizr'
import {OrgEntities} from 'src/types/schemas'

// Action Types
export const SET_ORGS = 'SET_ORGS'
export const SET_ORG = 'SET_ORG'
export const ADD_ORG = 'ADD_ORG'

export type Action =
  | ReturnType<typeof setOrgs>
  | ReturnType<typeof addOrg>
  | ReturnType<typeof setOrg>

export const setOrgs = (
  status: RemoteDataState,
  schema?: NormalizedSchema<OrgEntities, string[]>
) =>
  ({
    type: SET_ORGS,
    status,
    schema,
  } as const)

export const addOrg = (org: Organization) =>
  ({
    type: ADD_ORG,
    org,
  } as const)

export const setOrg = (org: Organization) =>
  ({
    type: SET_ORG,
    org,
  } as const)
