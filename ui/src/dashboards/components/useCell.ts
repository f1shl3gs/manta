// Libraries
import constate from 'constate';
import { useCallback, useEffect, useState } from 'react';

// Types
import { Cell } from 'types/Dashboard';

// Hooks
import { useParams } from 'react-router-dom';
import { CachePolicies, useFetch } from 'use-http';
import remoteDataState from 'utils/rds';

const [CellProvider, useCell] = constate(
  () => {
    const [cell, setCell] = useState<Cell>();

    const { cellID, dashboardID } = useParams<{ cellID: string, dashboardID: string }>();
    const { loading, error, patch, get } = useFetch(`/api/v1/dashboards/${dashboardID}/cells/${cellID}?a=b`, {
      cachePolicy: CachePolicies.NO_CACHE
    });

    useEffect(() => {
      get()
        .then(resp => {
          if (resp.viewProperties === undefined) {
            resp.viewProperties = {
              type: 'xy',
              xColumn: 'time',
              yColumn: 'value',
              axes: {
                x: {},
                y: {}
              },
              queries: [
                {
                  text: '',
                  hidden: false
                }
              ]
            };
          }

          console.log('set done', resp);
          setCell(resp);
        })
        .catch(err => {
          console.log('get failed ---- ', err);
        })
        .finally(() => {

        });
    }, []);

    const updateCell = useCallback((next: Cell) => {
      setCell(next);
      return patch(next);
    }, [cell]);

    console.log('rds', remoteDataState(cell, error, loading), 'cell', cell);

    return {
      cell,
      setCell,
      loading,
      error,
      updateCell,
      remoteDataState: remoteDataState(cell, error, loading)
    };
  },
  // useCell
  value => value
);

export {
  CellProvider,
  useCell
};