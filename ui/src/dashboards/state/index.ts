import constate from "constate";

import { Dashboards } from "types";
import { useOrgID } from "../../shared/state/organization/organization";
import { CachePolicies, useFetch } from "use-http";

const [DashboardsProvider, useDashboards] = constate(
  () => {
    const orgID = useOrgID();
    const { data, error, loading, get } = useFetch<Dashboards>(`/api/v1/dashboards?orgID=${orgID}`, {
      cache: CachePolicies.NO_CACHE
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