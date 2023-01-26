// Libraries
import {Dispatch} from 'react'
import {normalize} from 'normalizr'

// Types
import {RemoteDataState} from '@influxdata/clockface'
import {Organization} from 'src/types/organization'
import {OrgEntities} from 'src/types/schemas'
import {arrayOfOrgs, orgSchema} from 'src/schemas'

// Actions
import {Action, addOrg, setOrg, setOrgs} from 'src/organizations/actions'
import {
  info,
  notify,
  PublishNotificationAction,
} from 'src/shared/actions/notifications'
import {push, UpdateLocationActions} from '@lagunovsky/redux-react-router'

// Utils
import request from 'src/shared/utils/request'

// Constants
import {defaultErrorNotification} from 'src/shared/constants/notification'

export const getOrgs =
  () =>
  async (dispatch: Dispatch<Action>): Promise<Organization[]> => {
    try {
      const resp = await request('/api/v1/organizations')
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const orgs = resp.data

      dispatch(
        setOrgs(
          RemoteDataState.Done,
          normalize<Organization, OrgEntities, string[]>(orgs, arrayOfOrgs)
        )
      )

      return orgs
    } catch (err) {
      console.error(err)

      dispatch(setOrgs(RemoteDataState.Error, null))
    }
  }

export const createOrg =
  (org: Organization) =>
  async (
    dispatch: Dispatch<
      Action | PublishNotificationAction | UpdateLocationActions
    >
  ): Promise<void> => {
    try {
      const resp = await request('/api/v1/organizations', {
        method: 'POST',
        body: org,
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      dispatch(info('Create Organization success'))

      const norm = normalize<Organization, OrgEntities, string>(
        resp.data,
        orgSchema
      )

      dispatch(addOrg(norm))
      dispatch(setOrg(resp.data))
      dispatch(push(`/orgs/${norm.result}`))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Create Org failed, ${err}`,
        })
      )
    }
  }
