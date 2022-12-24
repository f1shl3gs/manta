// Libraries
import {normalize} from 'normalizr'
import {RemoteDataState} from '@influxdata/clockface'

// Types
import {GetState} from 'src/types/stores'
import {Check, CheckStatus} from 'src/types/checks'
import {CheckEntities} from 'src/types/schemas'
import {ResourceType} from 'src/types/resources'

// Actions
import {notifyError} from 'src/shared/actions/notifications'
import {removeCheck, setCheck, setChecks} from 'src/checks/actions/creators'

// Selectors
import {getOrg} from 'src/organizations/selectors'
import {getByID} from 'src/resources/selectors'

// Schemas
import {arrayOfChecks, checkSchema} from 'src/schemas/checks'

// Utils
import request from 'src/shared/utils/request'
import {setCheckBuilder} from './builder'
import {setTimeRange} from '../../shared/actions/timeRange'
import {pastHourTimeRange} from '../../shared/constants/timeRange'
import {back} from '@lagunovsky/redux-react-router'

export const createCheck =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const {id: orgID} = getOrg(state)

    const check = {
      name: 'name this check',
      desc: 'desc this check',
      status: 'active',
      cron: '@every 5s',
      orgID: orgID,
      query: 'rate(process_cpu_seconds_total[1m] * 100',
      conditions: [
        {
          status: 'warn',
          pending: '10m',
          threshold: {
            type: 'gt',
            value: 100,
          },
        },
      ],
    }

    try {
      const resp = await request(`/api/v1/checks`, {
        method: 'POST',
        body: check,
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Check, CheckEntities, string>(
        resp.data,
        checkSchema
      )

      dispatch(setCheck(norm.result, RemoteDataState.Done, norm))
    } catch (err) {
      console.log(err)

      dispatch(notifyError(`Create check failed, ${err}`))
    }
  }

export const getChecks =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const {id: orgID} = getOrg(state)

    try {
      const resp = await request(`/api/v1/checks?orgID=${orgID}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const checks = normalize<Check, CheckEntities, string[]>(
        resp.data,
        arrayOfChecks
      )

      dispatch(setChecks(RemoteDataState.Done, checks))
    } catch (err) {
      console.log(`get checks failed, ${err.message}`)

      dispatch(setChecks(RemoteDataState.Error))
      dispatch(notifyError(`get check failed, ${err}`))
    }
  }

export interface CheckUpdate {
  name?: string
  desc?: string
  activeStatus?: CheckStatus
}

export const patchCheck =
  (id: string, updates: CheckUpdate) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const current = getByID<Check>(state, ResourceType.Checks, id)
    const check = {
      ...current,
      ...updates,
    }

    try {
      const resp = await request(`/api/v1/checks/${id}`, {
        method: 'PATCH',
        body: {
          name: check.name,
          desc: check.desc,
          status: check.activeStatus,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const updated = normalize<Check, CheckEntities, string>(
        resp.data,
        checkSchema
      )
      dispatch(setCheck(id, RemoteDataState.Done, updated))
    } catch (err) {
      console.error(err)
      dispatch(notifyError(`Update check failed, ${err}`))
    }
  }

export const deleteCheck =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`/api/v1/checks/${id}`, {
        method: 'DELETE',
      })
      if (resp.status !== 204) {
        throw new Error(resp.data.message)
      }

      await dispatch(removeCheck(id))
    } catch (err) {
      console.log(err)

      dispatch(notifyError(`Delete check failed, ${err}`))
    }
  }

// builder
export const createCheckFromBuilder =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)
    const builder = state.checkBuilder

    try {
      const resp = await request(`/api/v1/checks`, {
        method: 'POST',
        body: {
          orgID: org.id,
          name: builder.name,
          desc: builder.desc,
          status: builder.activeStatus,
          cron: builder.cron,
          query: builder.query,
          conditions: Object.values(builder.conditions),
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Check, CheckEntities, string>(
        resp.data,
        checkSchema
      )
      await dispatch(setCheck(norm.result, RemoteDataState.Done, norm))
      dispatch(back())
    } catch (err) {
      console.error(err)

      dispatch(notifyError(`Create check failed, ${err}`))
    }
  }

export const getCheckForBuilder =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`/api/v1/checks/${id}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const check = normalize<Check, CheckEntities, string>(
        resp.data,
        checkSchema
      )

      await dispatch(setTimeRange(pastHourTimeRange))
      dispatch(setCheckBuilder(check.entities.checks[check.result]))
    } catch (err) {
      console.error(err)

      dispatch(notifyError(`Get check failed, ${err}`))
    }
  }
