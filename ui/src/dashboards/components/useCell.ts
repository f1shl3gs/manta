// Libraries
import constate from 'constate'
import {useHistory, useParams} from 'react-router-dom'
import {useCallback, useEffect, useState} from 'react'

// Hooks
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'shared/notification/useNotification'
import {useDashboard} from './useDashboard'

// Types
import {Cell} from 'types/Dashboard'
import {RemoteDataState} from '@influxdata/clockface'
import {Notification} from '../../types/Notification'

const cellUpdateFailed = (msg: string = 'unknown error'): Notification => ({
  ...defaultErrorNotification,
  message: `Failed to update cell: ${msg}`,
})

const cellFetchFailed = (msg: string = 'unknown error'): Notification => ({
  ...defaultErrorNotification,
  message: `Failed to fetch cell: ${msg}`,
})

const cellUpdateSuccess = (name: string): Notification => ({
  ...defaultSuccessNotification,
  message: `Update cell ${name} success`,
})

const [CellProvider, useCell] = constate(
  () => {
    const {reload} = useDashboard()
    const history = useHistory()
    const {cellID, dashboardID} = useParams<{
      cellID: string
      dashboardID: string
    }>()
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const {notify} = useNotification()
    const [cell, setCell] = useState<Cell>({
      desc: '',
      h: 0,
      id: '',
      maxW: 0,
      minH: 0,
      minW: 0,
      name: '',
      w: 0,
      x: 0,
      y: 0,
      viewProperties: {
        type: 'xy',
        xColumn: 'time',
        yColumn: 'value',
        hoverDimension: 'auto',
        geom: 'line',
        position: 'overlaid',
        axes: {
          x: {},
          y: {},
        },
        queries: [
          {
            text: '',
            hidden: false,
          },
        ],
      },
    })

    useEffect(() => {
      setLoading(RemoteDataState.Loading)

      fetch(`/api/v1/dashboards/${dashboardID}/cells/${cellID}`)
        .then(resp => resp.json())
        .then(data => {
          setCell(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setLoading(RemoteDataState.Error)
          notify(cellFetchFailed(err.message))
        })
    }, [])

    const updateCell = useCallback(
      (next: Cell) => {
        fetch(`/api/v1/dashboards/${dashboardID}/cells/${cellID}`, {
          method: 'PATCH',
          body: JSON.stringify(next),
        })
          .then(() => {
            notify(cellUpdateSuccess(next.name))
            history.goBack()
            reload()
          })
          .catch(err => {
            notify(cellUpdateFailed(err.message))
          })
      },
      [cellID, dashboardID, history, notify, reload]
    )

    const onRename = useCallback(
      (name: string) => {
        setCell({
          ...cell,
          name,
        })
      },
      [cell, setCell]
    )

    return {
      cell,
      setCell,
      loading,
      updateCell,
      onRename,
    }
  },
  // useCell
  value => value
)

export {CellProvider, useCell}
