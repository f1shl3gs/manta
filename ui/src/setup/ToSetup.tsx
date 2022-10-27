// Libraries
import {FC} from 'react'
import {Route, Routes, useNavigate} from 'react-router-dom'

// Components
import {
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import Setup from './SetupWizard'
import {LoginPage} from '../signin/LoginPage'

// Hooks
import useFetch from 'shared/useFetch'

export const ToSetup: FC = () => {
  const navigate = useNavigate()
  const {data = {allow: false}, loading} = useFetch('/api/v1/setup')

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <Routes>
        <Route path={'/setup'} element={<Setup />}></Route>
        <Route path={'/signin'} element={<LoginPage />}></Route>
      </Routes>
    </SpinnerContainer>
  )
}
