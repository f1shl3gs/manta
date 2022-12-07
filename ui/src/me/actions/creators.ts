import {RemoteDataState} from '@influxdata/clockface'

export const SET_ME = 'SET_ME'

export const setMe = (state: RemoteDataState, id?: string, name?: string) =>
  ({
    type: SET_ME,
    id,
    name,
    state,
  } as const)

export type Action = ReturnType<typeof setMe>
