import {RemoteDataState} from '@influxdata/clockface'
import {NormalizedSchema} from 'normalizr'
import {UserEntities} from 'src/types/user'

export const SET_MEMBERS = 'SET_MEMBERS'

//
export type UserSchema<R extends string | string[]> = NormalizedSchema<
  UserEntities,
  R
>

export type Action = ReturnType<typeof setMembers>

export const setMembers = (
  status: RemoteDataState,
  schema?: UserSchema<string[]>
) =>
  ({
    type: SET_MEMBERS,
    status,
    schema,
  } as const)
