import {Dispatch} from 'react'
import {normalize} from 'normalizr'

import {Action, addOrg, setOrgs} from 'src/organizations/actions'
import request from 'src/utils/request'
import {RemoteDataState} from '@influxdata/clockface'
import {Organization} from 'src/types/organization'

import {OrgEntities} from 'src/types/schemas'
import {arrayOfOrgs} from 'src/schemas'

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
  async (dispatch: Dispatch<Action>): Promise<void> => {
    try {
      const resp = await request('/api/v1/orgs', {method: 'POST', body: org})
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const created = resp.data

      dispatch(addOrg(created))
    } catch (err) {
      console.error(err)
    }
  }
