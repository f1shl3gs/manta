import constate from 'constate'
import {useFetch} from 'use-http'
import remoteDataState from '../utils/rds'

const [AuthenticationProvider, useAuth] = constate(() => {
  const {data, error, loading} = useFetch<boolean>('/api/v1/viewer', {}, [])

  return {
    data,
    loading: remoteDataState(data, error, loading),
  }
})

export {AuthenticationProvider, useAuth}
