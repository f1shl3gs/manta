import constate from 'constate'
import {useEffect, useState} from 'react'
import {RemoteDataState} from '@influxdata/clockface'

const [AuthenticationProvider, useAuth] = constate(() => {
  const [loading, setLoading] = useState(RemoteDataState.NotStarted)
  const [user, setUser] = useState({
    id: '',
    name: '',
  })

  useEffect(() => {
    setLoading(RemoteDataState.Loading)
    fetch('/api/v1/viewer')
      .then(resp => resp.json())
      .then(data => {
        setUser(data)
        setLoading(RemoteDataState.Done)
      })
      .catch(err => {
        console.log(err)
        setLoading(RemoteDataState.Error)
      })
  }, [])

  return {
    user,
    loading,
  }
})

export {AuthenticationProvider, useAuth}
