// Libraries
import {useCallback, useEffect, useState} from 'react'
import constate from 'constate'

// Hooks
import {useFetch} from '../../shared/useFetch'
import {useOrgID} from '../../shared/useOrg'
import {
  useNotification,
  defaultErrorNotification,
} from 'shared/notification/useNotification'

// Types
import {Variable} from 'types/Variable'
import {RemoteDataState} from '@influxdata/clockface'

const [VariablesProvider, useVariables] = constate(
  () => {
    const orgID = useOrgID()
    const {notify} = useNotification()
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const [variables, setVariables] = useState(new Array<Variable>())
    const {get} = useFetch<Variable[]>(`/api/v1/variables?orgID=${orgID}`)
    const {patch, del} = useFetch(`/api/v1/variables`)
    const [reload, setReload] = useState(0)

    useEffect(() => {
      setLoading(RemoteDataState.Loading)
      get()
        .then(list => {
          setVariables(list)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Fetch variables failed, ${err.message}`,
          })
        })
    }, [get, notify, reload])

    const onNameUpdate = useCallback(
      (id: string, name: string) => {
        patch(`/${id}`, {id, name})
          .then(() => setReload(prev => prev + 1))
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Update variable's name failed, err: ${err}`,
            })
          })
      },
      [notify, patch]
    )

    const onDescUpdate = useCallback(
      (id: string, desc: string) => {
        patch(`/${id}`, {id, desc})
          .then(() => setReload(prev => prev + 1))
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Update variable's desc failed, err: ${err} `,
            })
          })
      },
      [notify, patch]
    )

    const onDelete = useCallback(
      (v: Variable) => {
        del(`/${v.id}`)
          .then(() => setReload(prev => prev + 1))
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Delete Variable failed, err: ${v.name}`,
            })
          })
      },
      [del, notify]
    )

    return {
      loading,
      variables,
      onNameUpdate,
      onDescUpdate,
      onDelete,
    }
  },
  values => values
)

export {VariablesProvider, useVariables}
