// Libraries
import constate from 'constate';
import { useCallback, useEffect, useState } from 'react';

// Types
import { Cell } from 'types/Dashboard';

// Hooks
import { useParams } from 'react-router-dom';
import { CachePolicies, useFetch } from 'use-http';
import remoteDataState from '../../utils/rds';

const [CellProvider, useCell] = constate(
  () => {
    const [cell, setCell] = useState<Cell>();

    const { cellID, dashboardID } = useParams<{ cellID: string, dashboardID: string }>();
    const { data, loading, error, patch, get } = useFetch(`/api/v1/dashboards/${dashboardID}/cells/${cellID}?a=b`, {
      cachePolicy: CachePolicies.NO_CACHE
    });

    useEffect(() => {
      console.log('get')
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
                  hidden: false,
                }
              ]
            }
          }

          setCell(resp);
        })
        .catch(err => {
          console.log('get failed ---- ', err)
        })
        .finally(() => {
          console.log('get done')
        });
    }, []);

    const updateCell = useCallback((next: Cell) => {
      return patch(next);
    }, [cell]);

    console.log('cell', cell)

    return {
      cell,
      setCell,
      loading,
      error,
      updateCell,
      remoteDataState: remoteDataState(data, error, loading)
    };
  },
  // useCell
  value => value
);

export {
  CellProvider,
  useCell
};