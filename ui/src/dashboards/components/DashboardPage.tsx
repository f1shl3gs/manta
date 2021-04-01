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
import {TimeRangeProvider} from 'shared/useTimeRange'
import {AutoRefreshProvider} from 'shared/useAutoRefresh'
import {DashboardProvider, useDashboard} from './useDashboard'
import VariablesControlBar from './variablesControlBar/VariablesControlBar'
import {DndProvider} from 'react-dnd'
import {HTML5Backend} from 'react-dnd-html5-backend'

const dashRoute = `/orgs/:orgID/dashboards/:dashboardID`

const DashboardPage: React.FC = () => {
  const {cells, remoteDataState} = useDashboard()

  return (
    <TimeRangeProvider>
      <AutoRefreshProvider>
        <Page titleTag={'Dashboard'}>
          <SpinnerContainer
            loading={remoteDataState}
            spinnerComponent={<TechnoSpinner />}
          >
            <DashboardHeader />
            <VariablesControlBar />

            <Page.Contents
              fullWidth={true}
              scrollable={true}
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
      </AutoRefreshProvider>
    </TimeRangeProvider>
  )
}

const wrapper = () => (
  <DashboardProvider>
    <DashboardPage />
  </DashboardProvider>
)

export default wrapper
