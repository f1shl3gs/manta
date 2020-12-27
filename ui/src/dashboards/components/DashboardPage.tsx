import React from 'react';
import {
  Page,
  SpinnerContainer,
  TechnoSpinner
} from '@influxdata/clockface';
import { useFetch } from 'use-http';
import { Route, Switch, useParams } from 'react-router-dom';
import { useOrgID } from 'shared/hooks/useOrg';
import { Dashboard } from 'types/Dashboard';
import remoteDataState from '../../utils/rds';
import { DashboardProvider } from '../state/dashboard';
import Cells from './Cells';
import ViewEditorOverlay from './ViewEditorOverlay';
import DashboardHeader from './DashboardHeader';
import DashboardEmpty from './DashboardEmpty';
import { TimeRangeProvider } from '../../shared/useTimeRange';
import { AutoRefreshProvider } from '../../shared/useAutoRefresh';
import compose from '../../utils/compose';

const dashRoute = `/orgs/:orgID/dashboards/:dashboardID`;

const DashboardPage: React.FC = () => {
  const { dashboardID } = useParams<{ dashboardID: string }>();
  const orgID = useOrgID();
  const { data, error, loading } = useFetch<Dashboard>(
    `/api/v1/dashboards/${dashboardID}?orgID=${orgID}`,
    {},
    []
  );
  const rds = remoteDataState(loading, error);


  return (
    <AutoRefreshProvider>
      <TimeRangeProvider>
        <Page titleTag={'todo name'}>
          <DashboardHeader />

          <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
            <Page.Contents
              fullWidth={true}
              scrollable={true}
              className={'dashboard'}
            >
              {
                !!data?.cells ? (
                  <Cells />
                ) : (
                  <DashboardEmpty />
                )
              }
            </Page.Contents>
          </SpinnerContainer>
        </Page>

        <Switch>
          <Route
            path={`${dashRoute}/cells/:cellID/edit`}
            component={ViewEditorOverlay}
          />
        </Switch>
      </TimeRangeProvider>
    </AutoRefreshProvider>
  );
};

const wrapper = () => (
  <DashboardProvider>
    <DashboardPage />
  </DashboardProvider>
);

export default wrapper;
