import constate from "constate";
import { useParams } from "react-router-dom";
import { useFetch } from "use-http";

import { Dashboard } from "types";
import { Layout } from "react-grid-layout";

const [DashboardProvider, useCells, useReload] = constate(
  () => {
    const { dashboardID } = useParams<{ dashboardID: string }>();
    const { data, error, loading, get } = useFetch<Dashboard>(`/api/v1/dashboards/${dashboardID}`, {}, []);

    if (data !== undefined && data.cells === null) {
      data.cells = [];
    }

    // replace Cells
    const { put } = useFetch(`/api/v1/dashboards/${dashboardID}/cells`, {});
    const replaceCells = (layout: Layout[]) => {
      console.log('next', layout)
      
      const cells = layout.map(l => {
        return {
          id: l.i,
          x: l.x,
          y: l.y,
          w: l.w,
          h: l.h
        };
      });

      return put(cells)
    };

    return {
      loading,
      error,
      ...data,
      reload: get,
      replaceCells
    };
  },
  value => {
    return {
      cells: value.cells || [],
      setCells: value.replaceCells
    };
  },
  value => value.reload
);


export {
  DashboardProvider,
  useCells
};