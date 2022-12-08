// Libraries
import React, {FunctionComponent, useEffect, useState} from 'react'

// Components
import PageSpinner from 'src/shared/components/PageSpinner'

// Types
import {RemoteDataState} from '@influxdata/clockface'

// Utils
import request from 'src/shared/utils/request'
import {useDispatch} from 'react-redux';
import { push } from '@lagunovsky/redux-react-router'

interface Props {
  children: JSX.Element | JSX.Element[]
}

export const Setup: FunctionComponent<Props> = ({children}) => {
  const dispatch = useDispatch()
  const [loading, setLoading] = useState(RemoteDataState.Loading)

  useEffect(
    () => {
      request(`/api/v1/setup`)
        .then(resp => {
          if (resp.status !== 200) {
            throw new Error(resp.data.message)
          }

          const shouldSetup = resp.data?.allow || false
          if (shouldSetup) {
            // useNavigate looks like a better solution, but it is not
            //
            // https://github.com/remix-run/react-router/issues/7634
            dispatch(push(`/setup`))
          }

          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          console.error(err)

          setLoading(RemoteDataState.Error)
        })
    },
    [dispatch]
  )

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default Setup
