import constate from 'constate'
import {useOrgID} from '../shared/useOrg'
import {useCallback, useEffect, useState} from 'react'
import {
  defaultErrorNotification,
  useNotification,
} from '../shared/notification/useNotification'
import {RemoteDataState} from '@influxdata/clockface'

const [DashboardsProvider, useDashboards] = constate(
  () => {
    const orgID = useOrgID()
    const {notify} = useNotification()
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const [trigger, setTrigger] = useState(0)
    const [dashboards, setDashboards] = useState([])

    useEffect(() => {
      setLoading(RemoteDataState.Loading)
      fetch(`/api/v1/dashboards?orgID=${orgID}`)
        .then(resp => resp.json())
        .then(data => {
          setDashboards(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setLoading(RemoteDataState.Error)
          notify({
            ...defaultErrorNotification,
            message: 'Loading dashboards failed, err: ' + err.message,
          })
        })
    }, [notify, orgID, trigger])

    const reload = useCallback(() => {
      setTrigger(prevState => prevState + 1)
    }, [])

    const updateDashboard = useCallback(
      (id: string, upd: {name?: string; desc?: string}) => {
        fetch(`/api/v1/dashboards/${id}`, {
          method: 'PATCH',
          body: JSON.stringify(upd),
        })
          .then(() => {
            reload()
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: 'Update dashboard name failed',
            })
          })
      },
      [notify, reload]
    )

    return {
      loading,
      updateDashboard,
      dashboards,
      refresh: reload,
    }
  },
  // useDashboards
  value => value
)

export {DashboardsProvider, useDashboards}
