import React, {FunctionComponent} from 'react'
import EmptyQueryView from 'src/cells/components/EmptyQueryView'
import {ViewProperties} from 'src/types/cells'
import useQueryResult from 'src/shared/useQueryResult'
import View from 'src/visualization/View'
import {RemoteDataState} from '@influxdata/clockface'

interface Props {
  cellID?: string
  viewProperties: ViewProperties
}

const TimeSeries: FunctionComponent<Props> = props => {
  const {viewProperties} = props
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
