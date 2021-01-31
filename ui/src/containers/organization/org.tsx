import React from 'react'
import {Redirect, Route, Switch, useParams, withRouter} from 'react-router-dom'

import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import Otcl from 'otcls'

import {OrgProvider} from 'shared/useOrg'
import Todo from 'components/Todo'
import TracePage from 'traces'
import Logs from 'logs/Logs'
import DashboardsIndex from 'dashboards/dashboards'
import DashboardPage from 'dashboards/components/DashboardPage'
import Nav from 'layout/Nav'
import {useFetch} from 'use-http'
import remoteDataState from 'utils/rds'
import ProfilePage from '../../profile/ProfilePage'
import PluginsIndex from '../../plugins/PluginsIndex'
import PluginDetailsView from '../../plugins/PluginDetailsView'

const Org: React.FC = () => {
  const orgPath = '/orgs/:orgID'
  const {orgID} = useParams<{orgID: string}>()

  const {data, loading, error} = useFetch(`/api/v1/orgs/${orgID}`, {}, [])
  const rds = remoteDataState(data, error, loading)

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
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

          <Route path={`${orgPath}/otcls`} component={Otcl} />

          <Route path={`${orgPath}/traces`} component={TracePage} />
          <Route path={`${orgPath}/metrics`} component={Todo} />
          <Route path={`${orgPath}/logs`} component={Logs} />
          <Route path={`${orgPath}/alerting`} component={Todo} />
          <Route path={`${orgPath}/profile`} component={ProfilePage} />

          <Route
            path={`${orgPath}/dashboards/:dashboardID`}
            component={DashboardPage}
          />
          <Route path={`${orgPath}/dashboards`} component={DashboardsIndex} />
        </Switch>
      </OrgProvider>
    </SpinnerContainer>
  )
}

export default withRouter(Org)
