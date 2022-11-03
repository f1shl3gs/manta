import constate from 'constate'
import {useState} from 'react'
import {useNavigate} from 'react-router-dom'
import useFetch from 'shared/useFetch'

export const [OnboardProvider, useOnboard] = constate(() => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [organization, setOrganization] = useState('')
  const navigate = useNavigate()
  const {run: onboard} = useFetch(`/api/v1/setup`, {
    method: 'POST',
    body: {
      username,
      password,
      organization,
    },
    onSuccess: resp => {
      navigate(`/orgs/${resp?.org.id}`)
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
