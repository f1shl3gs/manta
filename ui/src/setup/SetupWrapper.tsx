import React, {FunctionComponent, useEffect} from 'react'
import useFetch from 'shared/useFetch'
import {RemoteDataState, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import {useNavigate} from 'react-router-dom'

interface Props {
  children: JSX.Element | JSX.Element[]
}

const SetupWrapper: FunctionComponent<Props> = ({children}) => {
  const {data, loading} = useFetch(`/api/v1/setup`)
  const navigate = useNavigate()

  useEffect(() => {
    if (loading === RemoteDataState.Done && data?.allow) {
      navigate(`/setup`)
    }
  }, [loading, data, navigate])

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner/>}>
      { children }
    </SpinnerContainer>
  )
}

export default SetupWrapper
