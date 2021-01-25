import React from 'react'
import {useHistory} from 'react-router-dom'
import {Provider} from 'use-http'
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import {AuthenticationProvider, useAuth} from '../shared/useAuthentication'

interface Props {
  children: React.ReactNode
}

const Authentication: React.FC<Props> = (props) => {
  const history = useHistory()
  const {children} = props
  const options = {
    interceptors: {
      // @ts-ignore
      response: async ({response}) => {
        if (response === undefined) {
          return undefined
        }

        if (response.status === 401) {
          history.push('/signin')
          return
        }

        return response
      },
    },
  }

  const {data, loading} = useAuth()

  console.log('data', data)

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      {children}
    </SpinnerContainer>
  )
}

const wrapped: React.FC = ({children}) => (
  <AuthenticationProvider>
    <Authentication>{children}</Authentication>
  </AuthenticationProvider>
)

export default wrapped
