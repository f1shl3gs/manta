// Libraries
import React, {FunctionComponent} from 'react'
import {RemoteDataState} from '@influxdata/clockface'

// Components
import EmptyQueryView from 'src/cells/components/EmptyQueryView'
import View from 'src/visualization/View'

// Types
import {ViewProperties} from 'src/types/cells'

// Hooks
import useQueryResult from 'src/shared/useQueryResult'

interface Props {
  viewProperties: ViewProperties
}

const TimeSeries: FunctionComponent<Props> = ({viewProperties}) => {
  const {result, loading, error} = useQueryResult(viewProperties.queries)

  return (
    <EmptyQueryView
      queries={viewProperties.queries}
      hasResults={loading === RemoteDataState.Done && result !== undefined}
      loading={loading}
      errorMessage={error}
    >
      <View result={result} properties={viewProperties} />
    </EmptyQueryView>
  )
}

export default TimeSeries
