// Libraries
import React from 'react'
import {Redirect, Route, Switch, useParams, withRouter} from 'react-router-dom'

// Components
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import {OrgProvider} from 'shared/useOrg'
import Todo from 'components/Todo'
import TracePage from 'traces'
import Logs from 'logs/Logs'
import Dashboards from 'dashboards/Dashboards'
import DashboardPage from 'dashboards/components/DashboardPage'
import Nav from 'layout/Nav'
import {Provider, useFetch} from 'shared/useFetch'
import remoteDataState from 'utils/rds'
import ProfilePage from '../../profile/ProfilePage'
import PluginsIndex from '../../plugins/PluginsIndex'
import PluginDetailsView from '../../plugins/PluginDetailsView'
import AlertsPage from '../../alerts/AlertsPage'
import Settings from '../../settings/Settings'
import Data from '../../data/Data'

const Org: React.FC = () => {
  const orgPath = '/orgs/:orgID'
  const {orgID} = useParams<{orgID: string}>()

  const {data, loading, error} = useFetch(`/api/v1/orgs/${orgID}`, {}, [])
  const rds = remoteDataState(data, error, loading)

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
      <Provider url={`/api/v1/orgs/${orgID}`}>
        <OrgProvider initialOrg={data}>
          <Nav />

          <Switch>
            {/* todo: memorize the path with localStorage? */}
            <Redirect exact from={`${orgPath}/`} to={`${orgPath}/dashboards`} />
            <Route exact path={`${orgPath}/plugins`} component={PluginsIndex} />
            <Route
              exact
              path={`${orgPath}/plugins/:id`}
              component={PluginDetailsView}
            />

            {/* OpenTelemetry Collectors */}
            {/* Data */}
            <Route path={`${orgPath}/data`} component={Data} />
            {/*<Route path={`${orgPath}/otcls`} component={Otcl} />*/}

            {/* Alerts */}
            <Route path={`${orgPath}/alerts`} component={AlertsPage} />

            {/* Traces */}
            <Route path={`${orgPath}/traces`} component={TracePage} />

            {/* Metrics */}
            <Route path={`${orgPath}/metrics`} component={Todo} />

            {/* Logs */}
            <Route path={`${orgPath}/logs`} component={Logs} />

            {/* Profile */}
            <Route path={`${orgPath}/profile`} component={ProfilePage} />

            {/* Dashboards */}
            <Route
              exact
              path={`${orgPath}/dashboards`}
              component={Dashboards}
            />
            <Route
              path={`${orgPath}/dashboards/:dashboardID`}
              component={DashboardPage}
            />

            {/* Settings */}
            <Route path={`${orgPath}/settings`} component={Settings} />
          </Switch>
        </OrgProvider>
      </Provider>
    </SpinnerContainer>
  )
}

export default withRouter(Org)
