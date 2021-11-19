import constate from 'constate'
import {useEffect, useState} from 'react'
import {useHistory, useLocation} from "react-router-dom";
import {RemoteDataState} from '@influxdata/clockface'
import {useFetch} from "./useFetch";

const [AuthenticationProvider, useAuth] = constate(() => {
  const history = useHistory();
  const location = useLocation();
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
        history.push(
            `/signin?returnTo=${encodeURIComponent(location.pathname)}`
        );
        console.log("get viewer info failed", err)
        setLoading(RemoteDataState.Error)
      })
  }, [])

  return {
    user,
    loading,
  }
})

export {AuthenticationProvider, useAuth}
