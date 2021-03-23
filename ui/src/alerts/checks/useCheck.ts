// Libraries
import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'

// types
import {Check, CheckStatusLevel} from '../../types/Check'
import {useFetch} from 'shared/useFetch'

// Utils
import {RemoteDataState} from '@influxdata/clockface'

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
}

const [CheckProvider, useCheck] = constate(
  (initialState: State) => {
    const [check, setCheck] = useState<Check>(defaultCheck)
    const [remoteDataState, setRemoteDataState] = useState(
      RemoteDataState.NotStarted
    )

    const {data, get, del, post, loading, error} = useFetch(
      `/api/v1/checks/${initialState.id}`,
      {}
    )

    useEffect(() => {
      setRemoteDataState(RemoteDataState.Loading)
      get()
        .then(data => {
          setCheck(data)
          setRemoteDataState(RemoteDataState.Done)
        })
        .catch(err => {
          setRemoteDataState(RemoteDataState.Error)
        })
    }, [])

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
