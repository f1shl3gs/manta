import React, {FunctionComponent} from 'react'
import useQueryResult from 'src/shared/useQueryResult'
import View from 'src/visualization/View'
import EmptyQueryView from 'src/cells/EmptyQueryView'
import {ViewProperties} from 'src/types/dashboard'

interface Props {
  viewProperties: ViewProperties
}

const TimeMachineVis: FunctionComponent<Props> = ({viewProperties}) => {
  const {result, loading} = useQueryResult(viewProperties.queries)

  return (
    <EmptyQueryView loading={loading} hasResults={result.table.length !== 0}>
      <View properties={viewProperties} result={result} />
    </EmptyQueryView>
  )
}

export default TimeMachineVis
