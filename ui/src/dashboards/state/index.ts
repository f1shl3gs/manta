import constate from "constate";
import { CachePolicies, useFetch } from "use-http";

import { Dashboards } from "types";
import { useOrgID } from "shared/state/organization/organization";

const [DashboardsProvider, useDashboards] = constate(
  () => {
    const orgID = useOrgID();
    const { data, error, loading, get } = useFetch<Dashboards>(`/api/v1/dashboards?orgID=${orgID}`, {
      cachePolicy: CachePolicies.NO_CACHE,
    }, []);

    return {
      error,
      loading,
      refresh: get,
      dashboards: data
    };
  }
);

export {
  DashboardsProvider,
  useDashboards
};