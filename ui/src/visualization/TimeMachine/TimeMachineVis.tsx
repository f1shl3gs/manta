import React, {FunctionComponent} from 'react'
import useQueryResult from 'src/shared/useQueryResult'
import View from 'src/visualization/View'
import EmptyQueryView from 'src/dashboards/components/Cell/EmptyQueryView'
import {useViewProperties} from 'src/visualization/TimeMachine/useViewProperties'

const TimeMachineVis: FunctionComponent = () => {
  const {viewProperties} = useViewProperties()
  const {result, loading} = useQueryResult(viewProperties.queries)

  return (
    <EmptyQueryView loading={loading} hasResults={result.table.length !== 0}>
      <View properties={viewProperties} result={result} />
    </EmptyQueryView>
  )
}

export default TimeMachineVis
