import React from 'react'
import {useHistory} from 'react-router-dom'
import {Provider} from 'use-http'

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

  return <Provider options={options}>{children}</Provider>
}

export default Authentication
