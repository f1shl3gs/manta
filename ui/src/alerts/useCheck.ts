// Libraries
import constate from 'constate'
import {useEffect, useState} from 'react'

// types
import {Check} from '../types/Check'
import {useFetch} from 'shared/useFetch'

// Utils
import remoteDataState from '../utils/rds'

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
    const {data, get, del, put, loading, error} = useFetch(
      `/api/v1/checks/${initialState.id}`,
      {}
    )
    useEffect(() => {
      get().then(data => {
        setCheck(data)
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

    return {
      check,
      remoteDataState: remoteDataState(data, error, loading),
      updateCheck,
    }
  },
  value => value
)

export {CheckProvider, useCheck}
