// Libraries
import constate from 'constate'
import {useEffect, useState} from 'react'

// types
import {Check} from '../types/Check'
import {useFetch} from 'use-http'

// Utils
import remoteDataState from '../utils/rds'

interface CheckUpdate {
  name: string
}

interface State {
  id: string
}

const [CheckProvider, useCheck] = constate(
  (initialState: State) => {
    const [check, setCheck] = useState<Check>()
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
