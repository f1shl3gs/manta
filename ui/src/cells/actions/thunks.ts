// Libraries
import {normalize} from 'normalizr'

// Types
import {Dashboard} from 'src/types/dashboards'
import {GetState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {DashboardEntities} from 'src/types/schemas'
import {dashboardSchema} from 'src/schemas'
import {cellSchema} from 'src/schemas/dashboards'
import {Cell, CellEntities, ViewProperties, ViewType} from 'src/types/cells'
import {RemoteDataState} from '@influxdata/clockface'

// Actions
import {notify} from 'src/shared/actions/notifications'
import {
  defaultDeletionNotification,
  defaultErrorNotification,
} from 'src/shared/constants/notification'
import * as creators from 'src/cells/actions/creators'
import {setCell} from 'src/cells/actions/creators'
import {back} from '@lagunovsky/redux-react-router'

// Selectors
import {getByID} from 'src/resources/selectors'

// Utils
import {get} from 'src/shared/utils/get'
import request from 'src/shared/utils/request'

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

const cellSize = (viewType: ViewType) => {
  if (viewType === 'gauge' || viewType === 'single-stat') {
    return {
      w: 3,
      h: 2,
    }
  } else {
    return {
      w: 4,
      h: 4,
    }
  }
}

export const createCell =
  (dashboardID: string, name: string) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const viewProperties = state.timeMachine.viewProperties
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
          ...cellSize(viewProperties.type),
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
      dispatch(back())
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
      if (resp.status !== 200) {
        throw new Error(resp.date.message)
      }

      const normCell = normalize<Cell, CellEntities, string>(
        resp.data,
        cellSchema
      )
      await dispatch(setCell(cellID, RemoteDataState.Done, normCell))
      dispatch(back())
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
