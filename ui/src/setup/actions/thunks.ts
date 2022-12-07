import {push} from '@lagunovsky/redux-react-router'
import {defaultErrorNotification} from 'src/shared/constants/notification'
import {notify} from 'src/shared/actions/notifications'
import {GetState} from 'src/types/stores'
import request from 'src/shared/utils/request'

const SETUP_URL = '/api/v1/setup'

export const initialUser =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const {username, password, organization} = state.setup

    try {
      const resp = await request(SETUP_URL, {
        method: 'POST',
        body: {
          username,
          password,
          organization,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const {id} = resp.data.org
      dispatch(push(`/orgs/${id}`))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Setup failed, ${err}`,
        })
      )
    }
  }
