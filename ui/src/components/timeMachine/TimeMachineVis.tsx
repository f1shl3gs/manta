// Libraries
import React from 'react'
import classnames from 'classnames'

// Components
import ErrorBoundary from 'shared/components/ErrorBoundary'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import ViewSwitcher from 'shared/components/ViewSwitcher'

// Hooks
import {useViewProperties} from 'shared/useViewProperties'
import {useQueries} from './useQueries'
import useQueryResult from 'dashboards/components/useQueryResult'

// Utils
import {RemoteDataState} from '@influxdata/clockface'

const TimeMachineVis: React.FC = () => {
  const {viewProperties} = useViewProperties()
  const timeMachineViewClassName = classnames('time-machine--view', {
    'time-machine--view__empty': false,
  })

  const {queries} = useQueries()
  // todo: handle errors
  const {result, errs} = useQueryResult(queries)

  return (
    <div className={timeMachineViewClassName}>
      <ErrorBoundary>
        {/*<ViewLoadingSpinner loading={rds} />*/}
        <EmptyQueryView
          loading={RemoteDataState.Done}
          queries={queries}
          hasResults={result?.table.length !== 0}
          errorMessage={errs}
        >
          <ViewSwitcher giraffeResult={result} properties={viewProperties} />
        </EmptyQueryView>
      </ErrorBoundary>
    </div>
  )
}

export default TimeMachineVis
