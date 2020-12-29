import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { CachePolicies, useFetch } from 'use-http';
import ViewSwitcher from './ViewSwitcher';
import moment from 'moment';
import { Config, fromRows, Plot, Table } from '@influxdata/giraffe';
import XYPlot from './XYPlot';
import { getFormatter } from '../../utils/vis';
import { useAutoRefresh } from '../../shared/useAutoRefresh';
import { AutoRefreshStatus } from '../../types/AutoRefresh';
import { useTimeRange } from '../../shared/useTimeRange';
import { TimeRange } from '../../types/TimeRanges';

interface Props {
  cellID?: string
}

type Result = {
  metric: {
    [key: string]: string
  }
  values: [
    [number, string]
  ]
}

const calculateRange = (timeRange: TimeRange) => {
  switch (timeRange.type) {
    case 'selectable-duration':
      const end = moment().unix();
      const start = end - timeRange.seconds;

      return {
        start,
        end,
        step: 14
      };

    default:
      const now = moment().unix();

      return {
        start: now - 3600,
        end: now,
        step: 14
      };
  }
};

const TimeSeries: React.FC<Props> = props => {
  const [table, setTable] = useState<Table>(() => fromRows([]));
  const { autoRefresh } = useAutoRefresh();
  const { timeRange } = useTimeRange();
  const url = `http://localhost:9090/api/v1/query_range`;
  const { get, loading } = useFetch(url, {
    cachePolicy: CachePolicies.NO_CACHE
  });

  const fetch = useCallback(() => {
    const { start, end, step } = calculateRange(timeRange);

    get(`?query=rate%28process_cpu_seconds_total%5B1m%5D%29+*+100&start=${start}&end=${end}&step=${step}`)
      .then(resp => {
        const rows = resp.data.result.map((item: Result) => {
          const { metric, values } = item;

          return values.map(val => {
            return {
              ...metric,
              time: val[0] * 1000,
              value: Number(val[1])
            };
          });
        }).flat();

        setTable(fromRows(rows));
      });
  }, [timeRange])

  useEffect(() => {
    fetch()
  }, [timeRange]);

  useEffect(() => {
    if (autoRefresh.status !== AutoRefreshStatus.Active) {
      return;
    }

    const timer = setInterval(fetch, autoRefresh.interval * 1000);

    return () => clearInterval(timer);
  }, [autoRefresh, fetch]);

  const xColumn = 'time';

  const xFormatter = getFormatter('time', {
    prefix: '',
    suffix: '',
    base: '10',
    timeZone: 'Local',
    timeFormat: 'YYYY/MM/DD HH:mm:ss'
  });

  const config = useMemo(() => {
    return {
      table,
      valueFormatters: {
        [xColumn]: xFormatter
      },
      layers: [
        {
          type: 'line',
          x: 'time',
          y: 'value',
          fill: ['job']
        }
      ]
    } as Config
  }, [table])

  if (loading) {
    return null
  }

  return (
    <Plot
      config={config}
    />
  );
};

export default TimeSeries;