import React from 'react'
import {useHistory, useLocation} from 'react-router-dom'
import {Provider} from 'shared/useFetch'
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import {useAuth} from '../shared/useAuthentication'

interface Props {
  children: React.ReactNode
}

const Authentication: React.FC<Props> = ({children}) => {
  const history = useHistory()
  const location = useLocation()

  const options = {
    interceptors: {
      // @ts-ignore
      response: async ({response}) => {
        if (response === undefined) {
          return undefined
        }

        if (response.status === 401) {
          history.push(
            `/signin?returnTo=${encodeURIComponent(location.pathname)}`
          )
          return
        }

        return response
      },
    },
  }

  const {loading} = useAuth()

  return (
    <Provider options={options}>
      <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
        {children}
      </SpinnerContainer>
    </Provider>
  )
}

export default Authentication
