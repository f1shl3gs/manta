// Libraries
import constate from 'constate'
import useFetch from 'src/shared/useFetch'
import {useNavigate} from 'react-router-dom'
import {useEffect} from 'react'
import {RemoteDataState} from '@influxdata/clockface'

interface User {
  id: string
  name: string
}

const [AuthenticationProvider, useAuthentication, useUser] = constate(
  () => {
    const {data, loading, error} = useFetch<User>('/api/v1/viewer')
    const navigate = useNavigate()

    useEffect(() => {
      switch (loading) {
        case RemoteDataState.Done:
          break
        case RemoteDataState.Error:
          navigate(
            `/signin?returnTo=${encodeURIComponent(window.location.pathname)}`
          )
          break
        default:
          break
      }
    }, [loading, error, data, navigate])

    return {
      user: data,
      loading,
    }
  },
  value => value,
  value => value.user || {id: '', name: ''}
)

export {AuthenticationProvider, useAuthentication, useUser}
