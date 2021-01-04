import constate from 'constate';
import { CachePolicies, useFetch } from 'use-http';
import { Dashboard } from '../types/Dashboard';
import { useOrgID } from '../shared/useOrg';
import remoteDataState from '../utils/rds';

const [DashboardsProvider, useDashboards] = constate(
  () => {
    const orgID = useOrgID();
    const { data, error, loading, get } = useFetch<Dashboard[]>(`/api/v1/dashboards?orgID=${orgID}`,
      {
        cachePolicy: CachePolicies.NO_CACHE
      },
      []);

    return {
      dashboards: data,
      remoteDataState: remoteDataState(data, error, loading),
      refresh: get
    };
  },
  // useDashboards
  value => value
);

export {
  DashboardsProvider,
  useDashboards
};