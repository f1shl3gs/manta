import request from 'src/utils/request'
import {getOrg} from 'src/organizations/selectors'
import {defaultErrorNotification} from 'src/constants/notification'
import {notify} from 'src/shared/actions/notifications'
import {normalize} from 'normalizr'
import {User, UserEntities} from 'src/types/user'
import {arrayOfMembers} from 'src/schemas/members'
import {setMembers} from 'src/members/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'
import {GetState} from 'src/types/stores'

export const getMembers =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)

    try {
      const resp = await request(`/api/v1/users?orgID=${org.id}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<User, UserEntities, string[]>(
        resp.data,
        arrayOfMembers
      )
      dispatch(setMembers(RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Get members failed, ${err}`,
        })
      )
    }
  }
