import React, {FunctionComponent, useEffect} from 'react'
import useFetch from 'shared/useFetch'
import {RemoteDataState} from '@influxdata/clockface'
import {Route, Routes, useNavigate} from 'react-router-dom'
import SetupWizard from './SetupWizard'
import PageSpinner from 'shared/components/PageSpinner'

interface SetupResp {
  allow: boolean
}

interface Props {
  children: JSX.Element | JSX.Element[]
}

const SetupWrapper: FunctionComponent<Props> = ({children}) => {
  const {data, loading} = useFetch<SetupResp>(`/api/v1/setup`)
  const navigate = useNavigate()
  const shouldSetup = data && data.allow

  useEffect(() => {
    if (loading === RemoteDataState.Done && data?.allow) {
      navigate(`/setup`)
    }
  }, [loading, data, navigate])

  return (
    <PageSpinner loading={loading}>
      {shouldSetup ? (
        <Routes>
          <Route path="/setup" element={<SetupWizard />} />
        </Routes>
      ) : (
        <>{children}</>
      )}
    </PageSpinner>
  )
}

export default SetupWrapper
