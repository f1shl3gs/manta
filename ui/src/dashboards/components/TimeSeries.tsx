import React, { useCallback, useEffect, useState } from 'react';

import { QueriesProvider } from '../../components/timeMachine/useQueries';
import { useViewProperties } from '../../shared/useViewProperties';
import ViewSwitcher from '../../shared/components/ViewSwitcher';
import { CachePolicies, useFetch } from 'use-http';
import { transformPromResp } from '../../utils/transform';
import { useAutoRefresh } from '../../shared/useAutoRefresh';
import { FromFluxResult, fromRows } from '@influxdata/giraffe';
import EmptyQueryView from '../../shared/components/EmptyQueryView';
import remoteDataState from '../../utils/rds';
import ErrorBoundary from '../../shared/components/ErrorBoundary';

interface Props {
  cellID?: string
}

const TimeSeries: React.FC<Props> = props => {
  const { viewProperties } = useViewProperties();
  const { start, end, step } = useAutoRefresh();
  const { queries } = viewProperties;
  const [errMsg, setErrMsg] = useState();
  const [result, setResult] = useState<Omit<FromFluxResult, 'schema'> | undefined>(() => {
    return {
      table: fromRows([]),
      fluxGroupKeyUnion: []
    };
  });

  const url = `http://localhost:9090/api/v1/query_range`;
  const { get, loading, error } = useFetch(url, {
    cachePolicy: CachePolicies.NO_CACHE,
    onError: () => {

    }
  });

  const fetch = useCallback(() => {
    get(`?query=${encodeURI(queries[0].text)}&start=${start}&end=${end}&step=${step}`)
      .then(resp => {
        if (resp.status !== 'success') {
          setErrMsg(resp.error);
          return;
        }

        setResult(transformPromResp(resp));
      });
  }, [start, end, step]);

  useEffect(() => {
    fetch();
  }, [start, end, step]);

  return (
    <ErrorBoundary>
      <QueriesProvider>
        {/*{result ? <ViewSwitcher properties={viewProperties} giraffeResult={result} /> : <EmptyQueryView/>}*/}
        <EmptyQueryView
          queries={queries}
          hasResults={result?.table.length !== 0}
          loading={remoteDataState(result, error, loading)}
          errorMessage={errMsg}
        >
          <ViewSwitcher properties={viewProperties} giraffeResult={result!} />
        </EmptyQueryView>
      </QueriesProvider>
    </ErrorBoundary>
  );
};

export default TimeSeries;