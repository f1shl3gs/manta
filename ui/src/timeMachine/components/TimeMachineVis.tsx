// Libraries
import React, {FunctionComponent, useEffect} from 'react'

// Components
import View from 'src/visualization/View'
import EmptyQueryView from 'src/cells/components/EmptyQueryView'

// Hooks
import {useDispatch, useSelector} from 'react-redux'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {loadView} from 'src/timeMachine/actions/thunks'

// Selectors
import {getTimeMachine} from 'src/timeMachine/selectors'

const TimeMachineVis: FunctionComponent = () => {
  const {
    queryResult: {state, result},
    viewProperties,
  } = useSelector((state: AppState) => {
    const timeMachine = getTimeMachine(state)

    return {
      queryResult: timeMachine.queryResult,
      viewProperties: timeMachine.viewProperties,
    }
  })

  const dispatch = useDispatch()
  useEffect(() => {
    dispatch(loadView())
  }, [dispatch])

  return (
    <EmptyQueryView loading={state} hasResults={result.table.length !== 0}>
      <View properties={viewProperties} result={result} />
    </EmptyQueryView>
  )
}

export default TimeMachineVis
