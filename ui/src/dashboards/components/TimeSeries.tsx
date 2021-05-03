// Libraries
import React from 'react'

// Components
import ViewSwitcher from 'shared/components/ViewSwitcher'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import {QueriesProvider} from 'components/timeMachine/useQueries'
import {useViewProperties} from 'shared/useViewProperties'

// Hooks
import useQueryResult from './useQueryResult'

// Types
import {RemoteDataState} from '@influxdata/clockface'

interface Props {
  cellID?: string
}

const TimeSeries: React.FC<Props> = () => {
  const {viewProperties} = useViewProperties()
  const {queries} = viewProperties
  const {result, errs} = useQueryResult(queries)

  return (
    <QueriesProvider>
      <EmptyQueryView
        queries={queries}
        hasResults={result !== undefined}
        loading={RemoteDataState.Done}
        errorMessage={errs}
      >
        <ViewSwitcher properties={viewProperties} giraffeResult={result!} />
      </EmptyQueryView>
    </QueriesProvider>
  )
}

export default TimeSeries
