// Libraries
import React, {FunctionComponent} from 'react'
// import {useNavigate} from 'react-router-dom'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useDispatch} from 'react-redux'
import {createCheck} from '../actions/thunks'

const CreateCheckButton: FunctionComponent = () => {
  // const navigate = useNavigate()
  const dispatch = useDispatch()

  const handleClick = () => {
    dispatch(createCheck())
  }

  return (
    <Button
      text={'Create Check'}
      testID={'create-check--button'}
      icon={IconFont.Plus_New}
      onClick={handleClick}
      color={ComponentColor.Primary}
    />
  )
}

export default CreateCheckButton
