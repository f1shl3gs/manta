// Libraries
import React, {FunctionComponent, useCallback} from 'react'

// Components
import {SquareButton, IconFont} from '@influxdata/clockface'
import {useDispatch} from 'react-redux'

// Actions
import {poll} from 'src/shared/actions/autoRefresh'

const AutoRefreshButton: FunctionComponent = () => {
  const dispatch = useDispatch()

  const handleClick = useCallback(() => {
    dispatch(poll())
  }, [dispatch])

  return <SquareButton icon={IconFont.Refresh_New} onClick={handleClick} />
}

export default AutoRefreshButton
