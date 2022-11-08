import React, {FunctionComponent, useEffect} from 'react'
import useFetch from 'src/shared/useFetch'
import {RemoteDataState} from '@influxdata/clockface'
import {Route, Routes, useNavigate} from 'react-router-dom'
import SetupWizard from 'src/setup/SetupWizard'
import PageSpinner from 'src/shared/components/PageSpinner'

interface SetupResp {
  allow: boolean
}

interface Props {
  children: JSX.Element | JSX.Element[]
}

const SetupWrapper: FunctionComponent<Props> = ({children}) => {
  const {data, loading} = useFetch<SetupResp>(`/api/v1/setup`)
  const navigate = useNavigate()
  const shouldSetup = data ? data.allow : false

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

// @ts-ignore
SetupWrapper.whyDidYouRender = true

export default SetupWrapper
