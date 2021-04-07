// Libraries
import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'

// types
import {Check, CheckStatusLevel} from '../../types/Check'
import {useFetch} from 'shared/useFetch'

// Utils
import {RemoteDataState} from '@influxdata/clockface'
import {
  defaultErrorNotification,
  useNotification,
} from '../../shared/notification/useNotification'

interface CheckUpdate {
  name: string
}

interface State {
  id: string
}

// todo: fix the status properties
const defaultCheck: Check = {
  lastRunStatus: '',
  latestCompleted: '',
  latestFailure: '',
  latestScheduled: '',
  latestSuccess: '',
  conditions: [],
  created: '',
  id: '',
  orgID: '',
  status: '',
  updated: '',
  name: '',
  desc: '',
  expr: '',
  labels: [],
}

const [CheckProvider, useCheck] = constate(
  (initialState: State) => {
    const [check, setCheck] = useState<Check>(defaultCheck)
    const {notify} = useNotification()
    const [remoteDataState, setRemoteDataState] = useState(
      RemoteDataState.NotStarted
    )

    const {get, post} = useFetch(`/api/v1/checks/${initialState.id}`, {})

    useEffect(() => {
      setRemoteDataState(RemoteDataState.Loading)
      get()
        .then(data => {
          // todo: dummy
          if (data.labels === null) {
            data.labels = []
          }

          setCheck(data)
          setRemoteDataState(RemoteDataState.Done)
        })
        .catch(err => {
          setRemoteDataState(RemoteDataState.Error)
          notify({
            ...defaultErrorNotification,
            message: `Fetch Check failed, err: ${err.message}`,
          })
        })
    }, [get, notify])

    const updateCheck = (udp: CheckUpdate) => {
      // @ts-ignore
      setCheck(prev => {
        return {
          ...prev,
          ...udp,
        }
      })
    }

    const onSave = useCallback(() => {
      return post(check)
    }, [check])

    const onRename = useCallback((name: string) => {
      setCheck(prev => {
        return {
          ...prev,
          name,
        }
      })
    }, [])

    const onAddCondition = useCallback((level: CheckStatusLevel) => {
      setCheck(prevState => {
        const {conditions} = prevState

        conditions.push({
          status: level,
          threshold: {
            type: 'gt',
            value: 0,
          },
        })
        return {
          ...prevState,
          conditions,
        }
      })
    }, [])

    return {
      ...check,
      onSave,
      onRename,
      onAddCondition,
      remoteDataState: remoteDataState,
      updateCheck,
    }
  },
  value => value
)

export {CheckProvider, useCheck}
