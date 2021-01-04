import React from 'react';
import { Route, Switch, useParams, withRouter } from 'react-router-dom';

import { SpinnerContainer, TechnoSpinner } from '@influxdata/clockface';
import Otcl from 'otcls';

import { OrgProvider } from 'shared/useOrg';
import Todo from 'components/Todo';
import TracePage from 'traces';
import remoteDataState from 'utils/rds';
import { useFetch } from 'use-http';
import Logs from 'logs/Logs';
import DashboardsIndex from '../../dashboards/dashboards';
import DashboardPage from '../../dashboards/components/DashboardPage';

const Org: React.FC = () => {
  const orgPath = '/orgs/:orgID';
  const { orgID } = useParams<{ orgID: string }>();

  const { data, loading, error } = useFetch(`/api/v1/orgs/${orgID}`, {}, []);
  const rds = remoteDataState(data, error, loading);

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
      <OrgProvider initialOrg={data}>
        <Switch>
          <Route path={`${orgPath}/otcls`} component={Otcl} />
          <Route path={`${orgPath}/traces`} component={TracePage} />
          <Route path={`${orgPath}/metrics`} component={Todo} />
          <Route path={`${orgPath}/logs`} component={Logs} />
          <Route path={`${orgPath}/alerting`} component={Todo} />

          <Route
            path={`${orgPath}/dashboards/:dashboardID`}
            component={DashboardPage}
          />
          <Route path={`${orgPath}/dashboards`} component={DashboardsIndex} />
        </Switch>
      </OrgProvider>
    </SpinnerContainer>
  );
};

export default withRouter(Org);
