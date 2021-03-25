// Libraries
import constate from 'constate'
import {useEffect, useState} from 'react'

// Types
import {NotificationEndpoint} from '../../client'
import {RemoteDataState} from '@influxdata/clockface'
import {useFetch} from '../../shared/useFetch'
import {useOrgID} from '../../shared/useOrg'

const [NotificationEndpointsProvider, useNotificationEndpoints] = constate(
  () => {
    const [endpoints, setEndpoints] = useState<NotificationEndpoint[]>([])
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const orgID = useOrgID()

    const {get} = useFetch(`/api/v1/notification_endpoints`)

    useEffect(() => {
      setLoading(RemoteDataState.Loading)

      get(`?orgID=${orgID}`)
        .then(data => {
          setEndpoints(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setLoading(RemoteDataState.Error)
        })
    }, [get, orgID])

    return {
      endpoints,
      loading,
      reload: get,
    }
  },
  value => value
)

export {NotificationEndpointsProvider, useNotificationEndpoints}
