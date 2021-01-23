// Libraries
import React from 'react'

// Components
import {
  AlignItems,
  AppWrapper,
  FunnelPage,
  InfluxDBCloudLogo,
  Panel,
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import SigninForm from './SigninForm'
import VersionInfo from './VersionInfo'

const Signin: React.FC = () => {
  return (
    <SpinnerContainer
      loading={RemoteDataState.Done}
      spinnerComponent={<TechnoSpinner />}
    >
      <AppWrapper>
        <FunnelPage className={'signin-page'}>
          <Panel className={'signin-page--panel'}>
            <Panel.Body alignItems={AlignItems.Center}>
              <div className={'signin-page--cubo'} />
              <InfluxDBCloudLogo
                cloud={false}
                className={'signin-page--logo'}
              />
              <SigninForm />
            </Panel.Body>

            <Panel.Footer>
              <VersionInfo />
            </Panel.Footer>
          </Panel>
        </FunnelPage>
      </AppWrapper>
    </SpinnerContainer>
  )
}

export default Signin
