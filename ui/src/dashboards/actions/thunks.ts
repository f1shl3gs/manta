import {normalize} from 'normalizr'
import {Dispatch} from 'react'

// Types
import {Layout} from 'react-grid-layout'

import request from 'src/utils/request'
import {GetState} from 'src/types/stores'
import {
  Action,
  editDashboard,
  removeDashboard,
  setDashboard,
  setDashboards,
} from 'src/dashboards/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'
import {getAll, getByID, getStatus} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {getOrg} from 'src/organizations/selectors'
import {DashboardEntities} from 'src/types/schemas'
import {Dashboard} from 'src/types/dashboards'
import {Cell} from 'src/types/cells'
import {arrayOfDashboards, dashboardSchema} from 'src/schemas'
import {
  notify,
  PublishNotificationAction,
} from 'src/shared/actions/notifications'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
} from 'src/constants/notification'
import {push} from '@lagunovsky/redux-react-router'
import {setCells} from 'src/cells/actions/creators'
import {arrayOfCells} from 'src/schemas/dashboards'
import {CellEntities} from 'src/types/cells'
import {getDashboard} from '../selectors'

export const getDashboards =
  () =>
  async (
    dispatch: Dispatch<Action>,
    getState: GetState
  ): Promise<Dashboard[]> => {
    try {
      const state = getState()
      if (
        getStatus(state, ResourceType.Dashboards) === RemoteDataState.NotStarted
      ) {
        dispatch(setDashboards(RemoteDataState.Loading))
      }

      const org = getOrg(state)

      const resp = await request(`/api/v1/dashboards?orgID=${org.id}`)
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const dashbaords = normalize<Dashboard, DashboardEntities, string[]>(
        resp.data,
        arrayOfDashboards
      )

      dispatch(setDashboards(RemoteDataState.Done, dashbaords))

      if (!dashbaords.result.length) {
        return
      }

      return resp.data
    } catch (err) {
      console.error(`get dashboard failed, ${err.message}`)
      dispatch(setDashboards(RemoteDataState.Error))

      throw err
    }
  }

export const createDashboard =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    try {
      const org = getOrg(getState())
      const dashboard = {
        name: '',
        cells: [],
        orgID: org.id,
      }

      const resp = await request(`/api/v1/dashboards`, {
        method: 'POST',
        body: dashboard,
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const normalized = normalize<Dashboard, DashboardEntities, string>(
        resp.data,
        dashboardSchema
      )

      await dispatch(
        setDashboard(resp.data.id, RemoteDataState.Done, normalized)
      )

      dispatch(push(`/orgs/${org.id}/dashboards/${resp.data.id}`))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Create Dashboard failed, ${err}`,
        })
      )
    }
  }

export const createDashboardFromJSON =
  (raw: Dashboard) =>
  async (dispatch, getState: GetState): Promise<void> => {
    try {
      const state = getState()
      const org = getOrg(state)

      // normalize
      const dashboard = {
        ...raw,
        orgID: org.id,
      }

      const resp = await request('/api/v1/dashboards', {
        method: 'POST',
        body: dashboard,
      })

      const normalized = normalize<Dashboard, DashboardEntities, string>(
        resp.data,
        dashboardSchema
      )

      await dispatch(
        setDashboard(resp.data.id, RemoteDataState.Done, normalized)
      )
      dispatch(
        notify({
          ...defaultSuccessNotification,
          message: `Import Dashboard success`,
        })
      )
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Import Dashboard failed, ${err}`,
        })
      )
    }
  }

export const cloneDashboard =
  (id: string, name: string) =>
  async (dispatch, getState: GetState): Promise<void> => {
    try {
      const state = getState()
      const org = getOrg(state)

      const getResp = await request(`/api/v1/dashboards/${id}`)
      if (getResp.status !== 200) {
        throw new Error(getResp.data.message)
      }

      const {entities, result} = normalize<
        Dashboard,
        DashboardEntities,
        string
      >(getResp.data, dashboardSchema)
      const dash: Dashboard = entities.dashboards[result]

      const postResp = await request(`/api/v1/dashboards`, {
        method: 'POST',
        body: dash,
      })
      if (postResp.status !== 201) {
        throw new Error(postResp.data.message)
      }

      const normilized = normalize<Dashboard, DashboardEntities, string>(
        postResp.data,
        dashboardSchema
      )
      await dispatch(
        setDashboard(postResp.data.id, RemoteDataState.Done, normilized)
      )

      dispatch(push(`/orgs/${org.id}/dashboards/${postResp.data.id}`))
    } catch (err) {
      console.log(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Clone Dashboard ${name} failed, ${err}`,
        })
      )
    }
  }

export const deleteDashboard =
  (id: string, name: string) =>
  async (dispatch: Dispatch<Action>): Promise<void> => {
    dispatch(removeDashboard(id))

    try {
      const resp = await request(`/api/v1/dashboards/${id}`, {method: 'DELETE'})
      if (resp.status !== 204) {
        throw new Error(resp.data.message)
      }
    } catch (err) {
      console.error(err)

      notify({
        ...defaultErrorNotification,
        message: `Delete Dashboard ${name} failed, ${err}`,
      })
    }
  }

export const updateDashboard =
  (id: string, updates: Partial<Dashboard>) =>
  async (
    dispatch: Dispatch<Action | PublishNotificationAction>,
    getState: GetState
  ): Promise<void> => {
    const state = getState()

    const current = getByID<Dashboard>(state, ResourceType.Dashboards, id)

    const dashboard = {
      ...current,
      ...updates,
    }

    try {
      const resp = await request(`/api/v1/dashboards/${id}`, {
        method: 'PATCH',
        body: {
          name: dashboard.name,
          desc: dashboard.desc,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const updatedDashboard = normalize<Dashboard, DashboardEntities, string>(
        resp.data,
        dashboardSchema
      )

      dispatch(editDashboard(updatedDashboard))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Update Dashboard ${current.name} failed, ${err}`,
        })
      )
    }
  }

export const updateCells =
  (dashboardID: string, cells: Cell[]) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`/api/v1/dashboards/${dashboardID}/cells`, {
        method: 'PUT',
        body: cells,
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const updatedCells = cells.map(c => ({...c, dashboardID}))
      const normalized = normalize<Dashboard, DashboardEntities, string[]>(
        updatedCells,
        arrayOfCells
      )

      dispatch(setCells(dashboardID, RemoteDataState.Done, normalized))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Failed to update cells, ${err}`,
        })
      )
    }
  }

export const updateLayout =
  (layouts: Layout[]) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const dashboard = getDashboard(state)
    const cells = getAll<Cell>(state, ResourceType.Cells)
    const newCell = layouts.map(layout => {
      const cell = cells.find(item => item.id == layout.id)

      return {
        ...cell,
        id: layout.id,
        x: layout.x,
        y: layout.y,
        w: layout.w,
        h: layout.h,
      }
    })

    try {
      const resp = await request(`/api/v1/dashboards/${dashboard.id}/cells`, {
        method: 'POST',
        body: newCell,
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const normalized = normalize<Cell, CellEntities, string[]>(
        resp.data,
        arrayOfCells
      )
      dispatch(setCells(dashboard.id, RemoteDataState.Done, normalized))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Update cells failed, ${err}`,
        })
      )
    }
  }
