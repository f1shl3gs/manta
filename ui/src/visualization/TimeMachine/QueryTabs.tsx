// Libraries
import React, {FunctionComponent} from 'react'

// Components
import QueryTab from 'src/visualization/TimeMachine/QueryTab'
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  SquareButton,
} from '@influxdata/clockface'

// Hooks
import {useQueries} from 'src/visualization/TimeMachine/useViewProperties'

const QueryTabs: FunctionComponent = () => {
  const {queries, addQuery} = useQueries()

  return (
    <div className={'time-machine-queries--tabs'}>
      {queries.map((query, index) => (
        <QueryTab key={index} index={index} query={query} />
      ))}

      <SquareButton
        className={'time-machine-queries--new'}
        icon={IconFont.Plus_New}
        size={ComponentSize.Small}
        color={ComponentColor.Default}
        onClick={addQuery}
      />
    </div>
  )
}

export default QueryTabs
