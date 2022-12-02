import constate from 'constate'
import {Dashboard} from 'src/types/dashboard'
import {useCallback, useState} from 'react'
import useFetch from 'src/shared/useFetch'
import {useParams} from 'react-router-dom'
import {Layout} from 'react-grid-layout'
import {
  defaultErrorNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'

interface State {
  dashboard: Dashboard
}

const [DashboardProvider, useDashboard] = constate((state: State) => {
  const {dashboardId} = useParams()
  const notify = useNotify()
  const [dashboard, setDashboard] = useState(state.dashboard)
  const {loading, run: reload} = useFetch(`/api/v1/dashboards/${dashboardId}`, {
    onSuccess: data => {
      setDashboard(data)
    },
  })

  const {run: replaceCells} = useFetch(
    `/api/v1/dashboards/${dashboardId}/cells`,
    {
      method: 'PUT',
      onError: err => {
        notify({
          ...defaultErrorNotification,
          message: `Update dashbaord ${dashboard.name} failed, ${err}`,
        })
      },
    }
  )

  const onLayoutChange = useCallback(
    (layouts: Layout[]) => {
      const cells = layouts.map(l => {
        const cell = dashboard.cells.find(item => item.id == l.i)

        return {
          ...cell,
          id: l.i,
          x: l.x,
          y: l.y,
          w: l.w,
          h: l.h,
        }
      })

      replaceCells(cells)
    },
    [dashboard, replaceCells]
  )

  const {run: patchDashboard} = useFetch(`/api/v1/dashboards/${dashboardId}`, {
    method: 'PATCH',
  })

  const onRename = useCallback(
    (name: string) => {
      patchDashboard({name})
    },
    [patchDashboard]
  )

  return {
    ...dashboard,
    reload,
    loading,
    onLayoutChange,
    onRename,
  }
})

export {DashboardProvider, useDashboard}
