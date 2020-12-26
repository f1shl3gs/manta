import React from 'react';
import {
  Page,
  SpinnerContainer,
  TechnoSpinner
} from '@influxdata/clockface';
import { useFetch } from 'use-http';
import { Route, Switch, useParams } from 'react-router-dom';
import { useOrgID } from '../../shared/state/organization/organization';
import { Dashboard } from '../../types';
import remoteDataState from '../../utils/rds';
import { AutoRefreshOption } from '../../types/AutoRefresh';
import { DashboardProvider } from '../state/dashboard';
import Cells from './Cells';
import ViewEditorOverlay from './ViewEditorOverlay';
import DashboardHeader from './DashboardHeader';
import DashboardEmpty from './DashboardEmpty';

const autoRefreshDropdownOptions: AutoRefreshOption[] = [
  {
    label: 'pause',
    seconds: 0
  },
  {
    label: 'Last 5m',
    seconds: 5 * 60
  },
  {
    label: 'Last 15m',
    seconds: 15 * 60
  },
  {
    label: 'Last 30m',
    seconds: 30 * 60
  },
  {
    label: 'Last 1h',
    seconds: 60 * 60
  },
  {
    label: 'Last 3h',
    seconds: 3 * 60 * 60
  },
  {
    label: 'Last 6h',
    seconds: 6 * 60 * 60
  }
];

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
    <>
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
    </>
  );
};

const wrapper = () => (
  <DashboardProvider>
    <DashboardPage />
  </DashboardProvider>
);

export default wrapper;
