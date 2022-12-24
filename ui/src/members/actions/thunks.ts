// Libraries
import {normalize} from 'normalizr'

// Utils
import request from 'src/shared/utils/request'

// Selector
import {getOrg} from 'src/organizations/selectors'

// Types
import {User, UserEntities} from 'src/types/user'
import {GetState} from 'src/types/stores'
import {RemoteDataState} from '@influxdata/clockface'

// Constants
import {defaultErrorNotification} from 'src/shared/constants/notification'

// Actions
import {notify} from 'src/shared/actions/notifications'
import {arrayOfMembers} from 'src/schemas/members'
import {setMembers} from 'src/members/actions/creators'

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
