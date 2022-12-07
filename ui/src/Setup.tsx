import React, {FunctionComponent, useEffect, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import request from './utils/request'
import {RemoteDataState} from '@influxdata/clockface'
import PageSpinner from './shared/components/PageSpinner'

interface Props {
  children: JSX.Element | JSX.Element[]
}

export const Setup: FunctionComponent<Props> = ({children}) => {
  const [loading, setLoading] = useState(RemoteDataState.Loading)
  const navigate = useNavigate()

  useEffect(
    () => {
      request(`/api/v1/setup`)
        .then(resp => {
          if (resp.status !== 200) {
            throw new Error(resp.data.message)
          }

          const shouldSetup = resp.data?.allow || false
          if (shouldSetup) {
            navigate('/setup')
          }

          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          console.error(err)

          setLoading(RemoteDataState.Error)
        })
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    []
  )

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default Setup
