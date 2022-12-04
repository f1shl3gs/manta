// Libraries
import React, {FunctionComponent, useCallback} from 'react'

// Components
import {SquareButton, IconFont} from '@influxdata/clockface'
import {useDispatch, useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'

// Actions
import {poll} from 'src/shared/actions/autoRefresh'

const AutoRefreshButton: FunctionComponent = () => {
  const dispatch = useDispatch()
  const timeRange = useSelector((state: AppState) => state.timeRange)

  const handleClick = useCallback(() => {
    dispatch(poll(timeRange))
  }, [dispatch, timeRange])

  return <SquareButton icon={IconFont.Refresh_New} onClick={handleClick} />
}

export default AutoRefreshButton
