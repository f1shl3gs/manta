import constate from 'constate'
import {useFetch} from 'use-http'
import {useHistory} from 'react-router-dom'

const [AuthenticationProvider, useAuth] = constate(() => {
  const {data, error, loading} = useFetch('/api/v1/viewer', {}, [])
})

export {AuthenticationProvider, useAuth}
