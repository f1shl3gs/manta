// Libraries
import React, {FunctionComponent} from 'react'

// Components
import QueryTab from 'src/timeMachine/QueryTab'
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  SquareButton,
} from '@influxdata/clockface'

// Hooks
import {useDispatch, useSelector} from 'react-redux'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {addQuery} from 'src/timeMachine/actions'

const QueryTabs: FunctionComponent = () => {
  const dispatch = useDispatch()
  const queries = useSelector((state: AppState) => {
    return state.timeMachine.queries
  })

  const handleClick = () => {
    dispatch(addQuery())
  }

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
        onClick={handleClick}
      />
    </div>
  )
}

export default QueryTabs
