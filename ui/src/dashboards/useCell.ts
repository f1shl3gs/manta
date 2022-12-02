import constate from 'constate'
import {useNavigate, useParams} from 'react-router-dom'
import {useCallback, useState} from 'react'

import {Cell, ViewProperties} from 'src/types/dashboard'
import useFetch from 'src/shared/useFetch'
import {
  defaultErrorNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'
import {defaultViewProperties} from 'src/constants/dashboard'

const defaultCell: Cell = {
  desc: '',
  h: 4,
  id: '',
  maxW: 0,
  minH: 0,
  minW: 0,
  name: '',
  w: 4,
  x: 0,
  y: 0,
  viewProperties: defaultViewProperties,
}

interface State {
  cell?: Cell
}

const [CellProvider, useCell] = constate((state: State) => {
  const {orgID, cellID, dashboardID} = useParams()
  const navigate = useNavigate()
  const notify = useNotify()

  // TODO: reload
  const reload = () => {
    console.log('realod')
  }
  const [cell, setCell] = useState<Cell>(() => {
    const cell = state.cell ?? defaultCell
    if (!cell.viewProperties) {
      cell.viewProperties = defaultViewProperties
    }

    return cell
  })

  const {run: patch} = useFetch(
    `/api/v1/dashboards/${dashboardID}/cells/${cellID}`,
    {
      method: 'PATCH',
      onError: err => {
        notify({
          ...defaultErrorNotification,
          message: `Update cell ${cell.name} failed, ${err}`,
        })
      },
      onSuccess: _ => {
        navigate(`/orgs/${orgID}/dashboards/${dashboardID}`)
      },
    }
  )

  const updateCell = useCallback(() => {
    console.log('update cell', cell)
    patch(cell)
  }, [cell, patch])

  const {run: del} = useFetch(
    `/api/v1/dashboards/${dashboardID}/cells/${cellID}`,
    {method: 'DELETE'}
  )
  const deleteCell = () => {
    del()
  }

  const {run: create} = useFetch(`/api/v1/dashboards/${dashboardID}/cells`, {
    method: 'POST',
    onSuccess: _ => {
      navigate(-1)
    },
  })
  const createCell = useCallback(() => {
    create(cell)
  }, [cell, create])

  const onRename = useCallback(
    name => {
      const newCell = {
        ...cell,
        name,
      }

      setCell(newCell)
      updateCell()
    },
    [cell, setCell, updateCell]
  )

  const setViewProperties = useCallback(
    (viewProperties: ViewProperties) => {
      setCell(prev => {
        return {
          ...prev,
          viewProperties,
        }
      })
    },
    [setCell]
  )

  return {
    cell,
    reload,
    createCell,
    updateCell,
    deleteCell,
    setCell,
    onRename,
    setViewProperties,
  }
})

export {CellProvider, useCell}
