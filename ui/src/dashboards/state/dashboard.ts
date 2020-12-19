import constate from "constate";
import { useParams } from "react-router-dom";
import { useFetch } from "use-http";
import { Dashboard } from "../../types";

const [DashboardProvider, useCells, useReload] = constate(
  () => {
    const { dashboardID } = useParams<{ dashboardID: string }>();
    const { data, error, loading, get } = useFetch<Dashboard>(`/api/v1/dashboards/${dashboardID}`, {}, []);

    if (data !== undefined && data.cells === null) {
      data.cells = [];
    }

    return {
      loading,
      error,
      ...data,
      reload: get
    };
  },
  value => {
    return value.cells || [];
  },
  value => value.reload
);


export {
  DashboardProvider,
  useCells
};