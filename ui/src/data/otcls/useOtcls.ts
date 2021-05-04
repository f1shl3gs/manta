// Libraries
import {useCallback, useEffect, useState} from 'react'

// Hooks
import {useOrgID} from '../../shared/useOrg'
import {
  defaultDeletionNotification,
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'shared/notification/useNotification'

// Types
import {RemoteDataState} from '@influxdata/clockface'
import {Otcl} from '../../types/otcl'
import constate from 'constate'

const [OtclsProvider, useOtcls] = constate(() => {
  const orgID = useOrgID()
  const [loading, setLoading] = useState(RemoteDataState.NotStarted)
  const [otcls, setOtcls] = useState<Otcl[]>([])
  const {notify} = useNotification()
  const [reload, setReload] = useState(0)

  useEffect(() => {
    setLoading(RemoteDataState.Loading)
    fetch(`/api/v1/otcls?orgID=${orgID}`)
      .then(resp => resp.json())
      .then(list => {
        setOtcls(list)
        setLoading(RemoteDataState.Done)
      })
      .catch(err => {
        notify({
          ...defaultErrorNotification,
          message: `Fetch otcls failed, ${err.message}`,
        })
        setLoading(RemoteDataState.Error)
      })
  }, [notify, orgID, reload])

  const onDelete = useCallback(
    (otcl: Otcl) => {
      fetch(`/api/v1/otcls/${otcl.id}`, {
        method: 'DELETE',
      })
        .then(() => {
          notify({
            ...defaultDeletionNotification,
            message: `Delete OTCL ${otcl.name} success`,
          })
          setReload(prev => prev + 1)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Delete OTCL ${otcl.name} failed, err: ${err.message}`,
          })
        })
    },
    [notify]
  )

  const onNameUpdate = useCallback(
    (id: string, name: string) => {
      fetch(`/api/v1/otcls/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          name,
        }),
      })
        .then(() => {
          notify({
            ...defaultSuccessNotification,
            message: `Update otcl name "${name}" success`,
          })
          setReload(prev => prev + 1)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Update otcl name failed, err: ${err.message}`,
          })
        })
    },
    [notify]
  )

  const onDescUpdate = useCallback(
    (id: string, desc: string) => {
      const otcl = otcls.find(item => item.id === id)
      if (!otcl) {
        return
      }

      fetch(`/api/v1/otcls/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          desc,
        }),
      })
        .then(() => {
          notify({
            ...defaultSuccessNotification,
            message: `Update ${otcl.name}'s desc success`,
          })
          setReload(prev => prev + 1)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Update ${otcl.name}'s desc failed`,
          })
        })
    },
    [notify, otcls]
  )

  return {
    loading,
    otcls,
    onDelete,
    onNameUpdate,
    onDescUpdate,
  }
})

export {OtclsProvider, useOtcls}
