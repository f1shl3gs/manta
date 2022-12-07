import {push} from '@lagunovsky/redux-react-router'

import request from 'src/shared/utils/request'
import {setMe} from './creators'
import {RemoteDataState} from '@influxdata/clockface'

export const getMe =
  () =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request('/api/v1/viewer')
      switch (resp.status) {
        case 401:
          dispatch(
            push(
              `/signin?returnTo=${encodeURIComponent(window.location.pathname)}`
            )
          )
          return

        case 200:
          const {id, name} = resp.data
          dispatch(setMe(RemoteDataState.Done, id, name))
          return

        default:
          dispatch(setMe(RemoteDataState.Error))
          throw new Error(resp.data.message)
      }
    } catch (err) {
      console.error(err)
    }
  }

export const signout = () => async (dispatch): Promise<void> => {
  try {
    const resp = await request(`/api/v1/signout`, {
      method: 'DELETE'
    })
    if (resp.status !== 200) {
      throw new Error(resp.date.message)
    }

    dispatch(push('/signin'))
  } catch (err) {
    console.error(err)
  }
}
