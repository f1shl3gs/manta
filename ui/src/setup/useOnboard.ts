import constate from 'constate'
import {useCallback, useState} from 'react'
import {useNavigate} from 'react-router-dom'

export const [OnboardProvider, useOnboard] = constate(() => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [organization, setOrganization] = useState('')
  const navigate = useNavigate()

  // TODO: error handle
  const onboard = useCallback(
    () =>
      fetch(`/api/v1/setup`, {
        method: 'POST',
        body: JSON.stringify({
          username,
          password,
          organization,
        }),
      })
        .then(data => data.json())
        .then(resp => {
          navigate(`/orgs/${resp.org.id}`)
        }),
    [username, password, organization]
  )

  return {
    username,
    password,
    organization,
    setUsername,
    setPassword,
    setOrganization,
    onboard,
  }
})
