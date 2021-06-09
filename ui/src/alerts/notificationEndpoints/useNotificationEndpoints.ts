// Libraries
import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'

// Hooks
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'shared/notification/useNotification'
import {useOrgID} from 'shared/useOrg'

// Types
import {NotificationEndpoint} from '../../client'
import {RemoteDataState} from '@influxdata/clockface'
import {useFetch} from '../../shared/useFetch'

const [NotificationEndpointsProvider, useNotificationEndpoints] = constate(
  () => {
    const [trigger, setTrigger] = useState(0)
    const [endpoints, setEndpoints] = useState<NotificationEndpoint[]>([])
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const orgID = useOrgID()
    const {notify} = useNotification()

    // const {data = [], error, loading } = useFetch(`/notification_endpoints`, {})

    useEffect(() => {
      setLoading(RemoteDataState.Loading)

      fetch(`/notification_endpoints?orgID=${orgID}`)
        .then(resp => resp.json())
        .then(data => {
          setEndpoints(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setLoading(RemoteDataState.Error)
          notify({
            ...defaultErrorNotification,
            message: `Fetch notification endpoints failed, err: ${err.message}`,
          })
        })
    }, [notify, orgID, trigger])

    const reload = useCallback(() => {
      setTrigger(prev => prev + 1)
    }, [])

    const patchNotificationEndpoint = useCallback(
      (id: string, upd: {name?: string; desc?: string}) => {
        fetch(`/api/v1/notification_endpoints/${id}`, {
          method: 'PATCH',
          body: JSON.stringify(upd),
        })
          .then(() => {
            notify({
              ...defaultSuccessNotification,
              message: `Update notification endpoint success`,
            })
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Update notification endpoint's ${
                upd.name ? 'name' : 'desc'
              } failed`,
            })
          })
      },
      [notify]
    )

    const deleteNotificationEndpoint = useCallback(
      (id: string) => {
        fetch(`/api/v1/notification_endpoints/${id}`, {
          method: 'DELETE',
        })
          .then(() => {
            reload()
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: 'Delete notification failed, err: ' + err.message,
            })
          })
      },
      [notify, reload]
    )

    return {
      loading,
      reload,
      endpoints,
      patchNotificationEndpoint,
      deleteNotificationEndpoint,
    }
  },
  value => value
)

export {NotificationEndpointsProvider, useNotificationEndpoints}
