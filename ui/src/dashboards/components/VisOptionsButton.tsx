import React, {useCallback} from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useDispatch, useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'
import {setViewingVisOptions} from 'src/timeMachine/actions'

const VisOptionsButton: React.FC = () => {
  const dispatch = useDispatch()
  const viewingVisOptions = useSelector((state: AppState) => {
    return state.timeMachine.viewingVisOptions
  })
  const handleClick = useCallback(() => {
    dispatch(setViewingVisOptions(!viewingVisOptions))
  }, [dispatch, viewingVisOptions])

  const color = viewingVisOptions
    ? ComponentColor.Primary
    : ComponentColor.Default

  return (
    <Button
      color={color}
      icon={IconFont.CogSolid_New}
      onClick={handleClick}
      text={'Customize'}
    />
  )
}

export default VisOptionsButton
