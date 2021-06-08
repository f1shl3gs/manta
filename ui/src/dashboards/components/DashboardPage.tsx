// Libraries
import React from 'react'
import {Route, Switch} from 'react-router-dom'

// Components
import {Page, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import ViewEditorOverlay from './ViewEditorOverlay'
import DashboardHeader from './DashboardHeader'
import DashboardEmpty from './DashboardEmpty'
import Cells from './Cells'

// Hooks
import {DashboardProvider, useDashboard} from './useDashboard'
import VariablesControlBar from './variablesControlBar/VariablesControlBar'
import {usePresentationMode} from '../../shared/usePresentationMode'

const dashRoute = `/orgs/:orgID/dashboards/:dashboardID`

const DashboardPage: React.FC = () => {
  const {inPresentationMode} = usePresentationMode()
  const {cells, remoteDataState, showVariablesControls} = useDashboard()

  return (
    <>
      <Page titleTag={'Dashboard'}>
        <SpinnerContainer
          loading={remoteDataState}
          spinnerComponent={<TechnoSpinner />}
        >
          <DashboardHeader />
          {/* todo: move VariablesControlBar to DashboardHeader */}
          {showVariablesControls && !inPresentationMode ? (
            <VariablesControlBar />
          ) : (
            <></>
          )}

          <Page.Contents
            fullWidth={true}
            scrollable={true}
            autoHideScrollbar={true}
            className={'dashboard'}
          >
            {cells.length > 0 ? <Cells /> : <DashboardEmpty />}
          </Page.Contents>
        </SpinnerContainer>
      </Page>

      <Switch>
        <Route
          path={`${dashRoute}/cells/:cellID/edit`}
          component={ViewEditorOverlay}
        />
      </Switch>
    </>
  )
}

const wrapper = () => (
  <DashboardProvider>
    <DashboardPage />
  </DashboardProvider>
)

export default wrapper
