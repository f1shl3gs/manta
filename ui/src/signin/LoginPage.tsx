// Libraries
import React, {FC} from 'react'

// Components
import {
  AlignItems,
  AppWrapper,
  FunnelPage,
  InfluxDBCloudLogo,
  Panel,
  PanelBody,
  PanelFooter,
} from '@influxdata/clockface'
import SignInForm from 'src/signin/SignInForm'
import {VersionInfo} from 'src/shared/components/VersionInfo'

const LoginPage: FC = () => {
  return (
    <AppWrapper>
      <FunnelPage className="signin-page" testID="signin-page">
        <Panel className="signin-page--panel">
          <PanelBody alignItems={AlignItems.Center}>
            <div className="signin-page--cubo" />

            <InfluxDBCloudLogo cloud={false} className="signin-page--logo" />

            <SignInForm />
          </PanelBody>

          <PanelFooter>
            <VersionInfo />
          </PanelFooter>
        </Panel>
      </FunnelPage>
    </AppWrapper>
  )
}

export default LoginPage
