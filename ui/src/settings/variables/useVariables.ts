// Libraries
import {useCallback, useEffect, useState} from 'react'
import constate from 'constate'

// Hooks
import {useFetch} from '../../shared/useFetch'
import {useOrgID} from '../../shared/useOrg'
import {
  useNotification,
  defaultErrorNotification,
} from '../../shared/notification/useNotification'

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
    }, [get, notify])

    const onNameUpdate = useCallback(
      (id: string, name: string) => {
        return patch(`${id}`, {id, name})
      },
      [patch]
    )

    const onDescUpdate = useCallback(
      (id: string, desc: string) => {
        return patch(`/${id}`, {id, desc})
      },
      [patch]
    )

    const onDelete = useCallback(
      (v: Variable) => {
        return del(`/${v.id}`)
      },
      [del]
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
