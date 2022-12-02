// Libraries
import React, {FunctionComponent} from 'react'

// Components
import View from 'src/visualization/View'
import EmptyQueryView from 'src/cells/components/EmptyQueryView'

// Hooks
import {useSelector} from 'react-redux'
import useQueryResult from 'src/shared/useQueryResult'

// Types
import {AppState} from 'src/types/stores'

const TimeMachineVis: FunctionComponent = () => {
  const {queries, viewProperties} = useSelector((state: AppState) => {
    return {
      queries: state.timeMachine.queries,
      viewProperties: state.timeMachine.viewProperties,
    }
  })

  const {result, loading} = useQueryResult(queries)

  return (
    <EmptyQueryView loading={loading} hasResults={result.table.length !== 0}>
      <View properties={viewProperties} result={result} />
    </EmptyQueryView>
  )
}

export default TimeMachineVis
