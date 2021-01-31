// Libraries
import React from 'react'
import classnames from 'classnames'

// Components
import ErrorBoundary from 'shared/components/ErrorBoundary'
import ViewLoadingSpinner from 'shared/components/ViewLoadingSpinner'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import ViewSwitcher from 'shared/components/ViewSwitcher'

// Hooks
import {useViewProperties} from 'shared/useViewProperties'
import {useAutoRefresh} from '../../shared/useAutoRefresh'
import {useQueries} from './useQueries'
import {useFetch} from 'use-http'

// Utils
import remoteDataState from 'utils/rds'
import {PromResp, transformPromResp} from 'utils/transform'

const TimeMachineVis: React.FC = () => {
  const {viewProperties} = useViewProperties()
  const timeMachineViewClassName = classnames('time-machine--view', {
    'time-machine--view__empty': false,
  })
  const {start, end, step} = useAutoRefresh()
  const {queries} = useQueries()

  const {data, loading, error} = useFetch<PromResp>(
    `http://localhost:9090/api/v1/query_range?query=rate%28process_cpu_seconds_total%5B1m%5D%29+*+100&start=${start}&end=${end}&step=${step}`,
    {},
    []
  )
  const rds = remoteDataState(data, error, loading)

  const gr = transformPromResp(data)

  return (
    <div className={timeMachineViewClassName}>
      <ErrorBoundary>
        <ViewLoadingSpinner loading={rds} />
        <EmptyQueryView
          loading={rds}
          queries={queries}
          hasResults={gr?.table.length !== 0}
          errorMessage={error !== undefined ? error.message : undefined}
        >
          <ViewSwitcher giraffeResult={gr!} properties={viewProperties} />
        </EmptyQueryView>
      </ErrorBoundary>
    </div>
  )
}

export default TimeMachineVis
