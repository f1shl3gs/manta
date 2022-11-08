import constate from 'constate'
import {useState} from 'react'
import useFetch from 'src/shared/useFetch'

export const [OnboardProvider, useOnboard] = constate(() => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [organization, setOrganization] = useState('')
  const {run: onboard} = useFetch(`/api/v1/setup`, {
    method: 'POST',
    body: {
      username,
      password,
      organization,
    },
    onSuccess: resp => {
      // /orgs/:orgID is not added to this Routes, so useNavigate
      // will not working
      window.location.href = `/orgs/${resp?.org.id}`
    },
  })

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
