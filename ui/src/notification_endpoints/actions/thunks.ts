// Libraries
import {normalize} from 'normalizr'

// Types
import {GetState} from 'src/types/stores'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'
import {NotificationEndpointEntities} from 'src/types/schemas'
import {RemoteDataState} from '@influxdata/clockface'

// Actions
import {
  removeNotificationEndpoint, setCurrentNotificationEndpoint,
  setNotificationEndpoint,
  setNotificationEndpoints,
} from 'src/notification_endpoints/actions/creators'
import {error, info} from 'src/shared/actions/notifications'
import {back} from '@lagunovsky/redux-react-router'

// Selectors
import {getOrgID} from 'src/organizations/selectors'
import {getEndpoint} from 'src/notification_endpoints/selectors'

// Utils
import request from 'src/shared/utils/request'

// Schema
import {
  arrayOfNotificationEndpoints,
  notificationEndpointSchema,
} from 'src/schemas/notificationEndpoints'

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

export const getNotificationEndpoint =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`${endpointPrefix}/${id}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<
        NotificationEndpoint,
        NotificationEndpointEntities,
        string
      >(resp.data, notificationEndpointSchema)
      dispatch(setNotificationEndpoint(id, RemoteDataState.Done, norm))
      dispatch(setCurrentNotificationEndpoint(resp.data))
    } catch (err) {
      console.error(err)

      dispatch(error(`Get notification endpoint failed, ${err}`))
    }
  }

export const createNotificationEndpoint =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const orgID = getOrgID(state)
    const endpoint = getEndpoint(state)

    try {
      const resp = await request(endpointPrefix, {
        method: 'POST',
        query: {
          orgID,
        },
        body: {
          ...endpoint,
          orgID,
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<
        NotificationEndpoint,
        NotificationEndpointEntities,
        string
      >(resp.data, notificationEndpointSchema)
      dispatch(setNotificationEndpoint(norm.result, RemoteDataState.Done, norm))
      dispatch(info('Create notification endpoint success'))
      dispatch(back())
    } catch (err) {
      console.error(err)

      dispatch(error(`Create notification endpoint failed, ${err}`))
    }
  }

export const deleteNotificationEndpoint =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`${endpointPrefix}/${id}`, {
        method: 'DELETE',
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      dispatch(removeNotificationEndpoint(id))
      dispatch(info(`Delete notification endpoint success`))
    } catch (err) {
      console.error(err)

      dispatch(error(`Delete notification endpoint failed`))
    }
  }

export interface NotificationEndpointUpdate {
  name?: string
  desc?: string
}

export const patchNotificationEndpoint =
  (id: string, upd: NotificationEndpointUpdate) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`${endpointPrefix}/${id}`, {
        method: 'PATCH',
        body: upd,
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<
        NotificationEndpoint,
        NotificationEndpointEntities,
        string
      >(resp.data, notificationEndpointSchema)
      dispatch(setNotificationEndpoint(id, RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(error(`Update notification endpoint failed, ${err}`))
    }
  }

export const updateNotificationEndpoint =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const endpoint = getEndpoint(state)

    try {
      const resp = await request(`${endpointPrefix}/${endpoint.id}`, {
        method: 'POST',
        body: endpoint,
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<
        NotificationEndpoint,
        NotificationEndpointEntities,
        string
      >(resp.data, notificationEndpointSchema)
      dispatch(setNotificationEndpoint(endpoint.id, RemoteDataState.Done, norm))
      dispatch(back())
    } catch (err) {
      console.error(err)

      dispatch(error(`Update notification endpoint failed, ${err}`))
    }
  }
