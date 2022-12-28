// Libraries
import {normalize} from 'normalizr'

// Types
import {Secret} from 'src/types/secrets'
import {SecretEntities} from 'src/types/schemas'
import {arrayOfSecrets, secretSchema} from 'src/schemas/secrets'
import {RemoteDataState} from '@influxdata/clockface'
import {GetState} from 'src/types/stores'

// Actions
import {removeSecret, setSecret, setSecrets} from 'src/secrets/actions/creators'
import {error} from 'src/shared/actions/notifications'
import {back} from '@lagunovsky/redux-react-router'

// Selectors
import {getOrgID} from 'src/organizations/selectors'

// Utils
import request from 'src/shared/utils/request'

interface CreateSecret {
  key: string
  value: string
}

export const upsertSecret =
  (newSecret: CreateSecret) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const orgID = getOrgID(state)

    try {
      const resp = await request(`/api/v1/secrets`, {
        method: 'POST',
        body: {
          ...newSecret,
          orgID,
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Secret, SecretEntities, string>(
        resp.data,
        secretSchema
      )

      dispatch(setSecret(RemoteDataState.Done, norm))
      dispatch(back())
    } catch (err) {
      console.error(err)

      dispatch(error(`create secret failed, ${err}`))
    }
  }

export const getSecrets =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const orgID = getOrgID(state)

    try {
      const resp = await request(`/api/v1/secrets?orgID=${orgID}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const secrets = normalize<Secret, SecretEntities, string[]>(
        resp.data,
        arrayOfSecrets
      )

      dispatch(setSecrets(RemoteDataState.Done, secrets))
    } catch (err) {
      console.error(err)

      dispatch(error(`Get secrets failed, ${err}`))
    }
  }

export const deleteSecret =
  (key: string) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const orgID = getOrgID(state)

    try {
      const resp = await request(`/api/v1/secrets/${key}`, {
        method: 'DELETE',
        query: {
          orgID,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      dispatch(removeSecret(key))
    } catch (err) {
      console.error(err)

      dispatch(error(`Delete secret ${key} failed, ${err}`))
    }
  }
