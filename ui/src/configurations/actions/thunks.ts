import {getOrg} from 'src/organizations/selectors'
import {notify} from 'src/shared/actions/notifications'
import {
  defaultDeletionNotification,
  defaultErrorNotification,
} from 'src/shared/constants/notification'
import {GetState} from 'src/types/stores'
import request from 'src/shared/utils/request'
import {getByID} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {Configuration} from 'src/types/configuration'
import {normalize} from 'normalizr'
import {ConfigurationEntities} from 'src/types/schemas'
import {arrayOfConfigurations, configurationSchema} from 'src/schemas'
import {
  removeConfig,
  setConfig,
  setConfigs,
} from 'src/configurations/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'

export const getConfigs =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)

    try {
      const resp = await request(`/api/v1/configurations?orgID=${org.id}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Configuration, ConfigurationEntities, string[]>(
        resp.data,
        arrayOfConfigurations
      )
      dispatch(setConfigs(RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Get configs failed, ${err}`,
        })
      )
    }
  }

export const createConfig =
  (name: string, desc: string, content: string) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)

    try {
      const resp = await request('/api/v1/configurations', {
        method: 'POST',
        body: {
          name,
          desc,
          data: content,
          orgID: org.id,
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Configuration, ConfigurationEntities, string>(
        resp.data,
        configurationSchema
      )

      dispatch(setConfig(norm.result, RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Create Configuration failed, ${err}`,
        })
      )
    }
  }

export const deleteConfig =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`/api/v1/configurations/${id}`, {
        method: 'DELETE',
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      dispatch(removeConfig(id))
      dispatch(
        notify({
          ...defaultDeletionNotification,
          message: 'Delete config success',
        })
      )
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Delete config failed, ${err}`,
        })
      )
    }
  }

export interface ConfigUpdate {
  name?: string
  desc?: string
  content?: string
}

export const updateConfig =
  (id: string, updates: ConfigUpdate) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const current = getByID<Configuration>(
      state,
      ResourceType.Configurations,
      id
    )

    const config = {
      ...current,
      ...updates,
    }

    try {
      const resp = await request(`/api/v1/configurations/${id}`, {
        method: 'PATCH',
        body: updates,
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Configuration, ConfigurationEntities, string>(
        resp.data,
        configurationSchema
      )

      dispatch(setConfig(norm.result, RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Update config ${config.name} failed, ${err}`,
        })
      )
    }
  }