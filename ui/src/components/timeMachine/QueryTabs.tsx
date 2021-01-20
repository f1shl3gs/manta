// Libraries
import React from 'react'

// Components
import QueryTab from './QueryTab'
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  SquareButton,
} from '@influxdata/clockface'

// Hooks
import {useQueries} from './useQueries'

const QueryTabs: React.FC = () => {
  const {queries, setActiveIndex, addQuery} = useQueries()

  return (
    <div className={'time-machine-queries--tabs'}>
      {queries.map((query, index) => (
        <QueryTab key={index} queryIndex={index} query={query} />
      ))}

      <SquareButton
        className={'time-machine-queries--new'}
        icon={IconFont.PlusSkinny}
        size={ComponentSize.Small}
        color={ComponentColor.Default}
        onClick={addQuery}
      />
    </div>
  )
}

export default QueryTabs
