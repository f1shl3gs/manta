import constate from 'constate';
import { useParams } from 'react-router-dom';
import { CachePolicies, useFetch } from 'use-http';

import { Dashboard } from 'types/Dashboard';
import { Layout } from 'react-grid-layout';
import { useCallback, useEffect, useState } from 'react';

const [DashboardProvider, useCells, useReload] = constate(
  () => {
    const { dashboardID } = useParams<{ dashboardID: string }>();
    const { error, loading, get } = useFetch<Dashboard>(
      `/api/v1/dashboards/${dashboardID}`,
      {
        cachePolicy: CachePolicies.NO_CACHE
      });

    const [dash, setDash] = useState<Dashboard>({
      id: '',
      created: '',
      updated: '',
      name: '',
      desc: '',
      orgID: '',
      cells: []
    });

    useEffect(() => {
      get()
        .then(res => setDash(res));
    }, []);

    // replace Cells
    const { put } = useFetch(`/api/v1/dashboards/${dashboardID}/cells`, {});
    const replaceCells = useCallback((layout: Layout[]) => {
      const cells = layout.map((l) => {
        return {
          id: l.i,
          x: l.x,
          y: l.y,
          w: l.w,
          h: l.h
        };
      });

      return put(cells);
    }, []);

    return {
      loading,
      error,
      ...dash,
      reload: get,
      replaceCells
    };
  },
  (value) => {
    return {
      cells: value.cells || [],
      setCells: value.replaceCells
    };
  },
  (value) => value.reload
);

export { DashboardProvider, useCells };
