import React from 'react';
import {
  Route,
  RouteComponentProps,
  Switch,
  useParams,
  withRouter,
} from 'react-router';

import {
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface';
import Otcl from 'otcl';

import { OrgProvider } from 'shared/state/organization/organization';
import useFetch from 'use-http';
import Todo from '../../components/Todo';

type Props = RouteComponentProps<{ orgID: string }>;

const Org: React.FC<Props> = (props) => {
  const orgPath = '/orgs/:orgID';
  const { orgID } = useParams();

  const { data, loading, error } = useFetch(`/api/v1/orgs/${orgID}`, {}, []);
  const rds = loading
    ? RemoteDataState.Loading
    : error !== undefined
    ? RemoteDataState.Error
    : RemoteDataState.Done;

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
      <OrgProvider initialOrg={data}>
        <Switch>
          <Route path={`${orgPath}/otcls`} component={Otcl} />
          <Route path={`${orgPath}/metrics`} component={Todo} />
          <Route path={`${orgPath}/logs`} component={Todo} />
          <Route path={`${orgPath}/traces`} component={Todo} />
          <Route path={`${orgPath}/alerting`} component={Todo} />
        </Switch>
      </OrgProvider>
    </SpinnerContainer>
  );
};

export default withRouter(Org);
