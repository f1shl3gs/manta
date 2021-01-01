// Libraries
import React from 'react';
import classnames from 'classnames';

// Components
import ErrorBoundary from 'shared/components/ErrorBoundary';
import ViewLoadingSpinner from 'shared/components/ViewLoadingSpinner';
import EmptyQueryView from 'shared/components/EmptyQueryView';
import ViewSwitcher from 'shared/components/ViewSwitcher';

// Types
import { ViewProperties } from 'types/Dashboard';
import { useFetch } from 'use-http';
import remoteDataState from '../../utils/rds';
import { PromResp, transformPromResp } from '../../utils/transform';

interface Props {
  viewProperties: ViewProperties
}

const TimeMachineVis: React.FC<Props> = props => {
  const {
    viewProperties
  } = props;

  const timeMachineViewClassName = classnames('time-machine--view', {
    'time-machine--view__empty': false
  });

  const {
    data,
    loading,
    error
  } = useFetch<PromResp>(`http://localhost:9090/api/v1/query_range?query=rate%28process_cpu_seconds_total%5B1m%5D%29+*+100&start=1609486814&end=1609490414&step=14`, {}, []);
  const rds = remoteDataState(loading, error);

  const gr = transformPromResp(data)

  return (
    <div className={timeMachineViewClassName}>
      <ErrorBoundary>
        <ViewLoadingSpinner loading={rds} />
        <EmptyQueryView
          loading={rds}
          queries={viewProperties.queries}
          // error={''}
        >
          <ViewSwitcher
            giraffeResult={gr}
            properties={viewProperties}
          />
        </EmptyQueryView>
      </ErrorBoundary>
    </div>
  );
};

export default TimeMachineVis;
