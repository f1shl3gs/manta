// Libraries
import {normalize} from 'normalizr'

// Types
import {GetState} from 'src/types/stores'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'
import {NotificationEndpointEntities} from 'src/types/schemas'
import {RemoteDataState} from '@influxdata/clockface'

// Selectors
import {getOrgID} from 'src/organizations/selectors'

// Utils
import request from 'src/shared/utils/request'

// Schema
import {arrayOfNotificationEndpoints} from 'src/schemas/notificationEndpoints'

// Actions
import {setNotificationEndpoints} from 'src/notification_endpoints/actions/creators'
import {error} from 'src/shared/actions/notifications'

const endpointPrefix = '/api/v1/notificationEndpoints'

export const getNotificationEndpoints =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const orgID = getOrgID(state)

    try {
      const resp = await request(endpointPrefix, {
        query: {
          orgID,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const eps = normalize<
        NotificationEndpoint,
        NotificationEndpointEntities,
        string[]
      >(resp.data, arrayOfNotificationEndpoints)

      dispatch(setNotificationEndpoints(RemoteDataState.Done, eps))
    } catch (err) {
      console.error(err)

      dispatch(error(`Get notification endpoints failed, ${err}`))
    }
  }
