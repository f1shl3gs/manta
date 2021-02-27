import {useCallback} from 'react'
import constate from 'constate'
import {useParams} from 'react-router-dom'
import {CachePolicies, useFetch} from 'shared/useFetch'

import remoteDataState from '../../utils/rds'
import {Cells, Dashboard} from '../../types/Dashboard'
import {Layout} from 'react-grid-layout'

const [DashboardProvider, useDashboard] = constate(
  () => {
    const {dashboardID} = useParams<{dashboardID: string}>()
    const {data, loading, error, get} = useFetch<Dashboard>(
      `/api/v1/dashboards/${dashboardID}`,
      {
        cachePolicy: CachePolicies.NO_CACHE,
      },
      []
    )

    const {post: update} = useFetch(`/api/v1/dashboards/${dashboardID}`, {})

    // onRename
    const onRename = useCallback((name: string) => {
      return update({
        name,
      })
    }, [])

    // addCell
    const {post: addCellPost} = useFetch(
      `/api/v1/dashboards/${dashboardID}/cells`,
      {}
    )
    const addCell = useCallback(() => {
      return addCellPost({
        w: 4,
        h: 4,
        x: 0,
        y: 0,
      }).then(() => {
        // refresh
        get()
      })
    }, [get])

    // delete cell by id
    const {del} = useFetch(`/api/v1/dashboards/${dashboardID}/cells/`, {})
    const deleteCell = useCallback(
      (cellID: string) => {
        return del(cellID)
      },
      [del]
    )

    // resetCells
    const {put} = useFetch(`/api/v1/dashboards/${dashboardID}/cells`, {})
    const onLayoutChange = useCallback(
      (layouts: Layout[]) => {
        const cells = layouts.map(l => {
          const cell = data?.cells.find(item => item.id == l.i)

          return {
            ...cell,
            id: l.i,
            x: l.x,
            y: l.y,
            w: l.w,
            h: l.h,
          }
        })

        return put(cells)
      },
      [data]
    )

    return {
      ...data,
      addCell,
      update,
      onRename,
      onLayoutChange,
      deleteCell,
      reload: get,
      remoteDataState: remoteDataState(data, error, loading),
    }
  },
  // useDashboard
  value => {
    return {
      ...value,
      cells: value.cells || ([] as Cells),
    }
  }
)

export {DashboardProvider, useDashboard}
