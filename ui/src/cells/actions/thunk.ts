import {get} from 'lodash'
import {notify} from 'src/shared/actions/notifications'
import {
  defaultDeletionNotification,
  defaultErrorNotification,
} from 'src/constants/notification'
import {Dashboard} from 'src/types/dashboards'
import {GetState} from 'src/types/stores'
import request from 'src/utils/request'
import * as creators from 'src/cells/actions/creators'
import {setCell} from 'src/cells/actions/creators'
import {getByID} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {normalize} from 'normalizr'
import {DashboardEntities} from 'src/types/schemas'
import {dashboardSchema} from 'src/schemas'
import {cellSchema} from 'src/schemas/dashboards'
import {Cell, CellEntities, ViewProperties} from 'src/types/cells'
import {RemoteDataState} from '@influxdata/clockface'

export const getCell =
  (cellID: string) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const dashboardID = get(state, 'resources.dashboards.current', '')

    try {
      const resp = await request(
        `/api/v1/dashboards/${dashboardID}/cells/${cellID}`
      )
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Cell, CellEntities, string>(resp.data, cellSchema)
      dispatch(setCell(cellID, RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(setCell(cellID, RemoteDataState.Error))
      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Get Cell failed, ${err}`,
        })
      )
    }
  }

export const createCell =
  (dashboardID: string, name: string, viewProperties: ViewProperties) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    let dashboard = getByID<Dashboard>(
      state,
      ResourceType.Dashboards,
      dashboardID
    )

    try {
      if (!dashboard) {
        // not exist or something wrong
        const resp = await request(`/api/v1/dashboards/${dashboardID}`)
        if (resp.status !== 200) {
          throw new Error(resp.data.message)
        }

        const normalized = normalize<Dashboard, DashboardEntities, string>(
          resp.data,
          dashboardSchema
        )

        const {entities, result} = normalized
        dashboard = entities[result]
      }

      // Create the cell
      const resp = await request(`/api/v1/dashboards/${dashboardID}/cells`, {
        method: 'POST',
        body: {
          name,
          viewProperties,
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const cellID = resp.data.id
      const normCell = normalize<Cell, CellEntities, string>(
        {
          ...resp.data,
          dashboardID,
        },
        cellSchema
      )

      dispatch(setCell(cellID, RemoteDataState.Done, normCell))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Create cell failed, ${err}`,
        })
      )
    }
  }

export const removeCell =
  (dashboardID: string, cellID: string) =>
  async (dispatch, _getState: GetState): Promise<void> => {
    try {
      const resp = await request(
        `/api/v1/dashboards/${dashboardID}/cells/${cellID}`,
        {
          method: 'DELETE',
        }
      )
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      dispatch(creators.removeCell(dashboardID, cellID))
      dispatch(
        notify({
          ...defaultDeletionNotification,
          message: `Cell deleted from dashboard`,
        })
      )
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Delete cell failed, ${err}`,
        })
      )
    }
  }

export interface CellUpdate {
  name?: string
  viewProperties?: ViewProperties
}

export const updateCell =
  (dashboardID: string, cellID: string, upd: CellUpdate) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const prev = getByID<Cell>(state, ResourceType.Cells, cellID)

    try {
      const resp = await request(
        `/api/v1/dashboards/${dashboardID}/cells/${cellID}`,
        {
          method: 'PATCH',
          body: {
            ...prev,
            ...upd,
          },
        }
      )
      if (resp.status !== 201) {
        throw new Error(resp.date.message)
      }
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Update cell failed, ${err}`,
        })
      )
    }
  }
