// Libraries
import React from 'react'
import {useDispatch} from 'react-redux'

// Components
import {Button, ComponentColor} from '@influxdata/clockface'

// Actions
import {loadView} from 'src/timeMachine/actions/thunks'

const SubmitQueryButton: React.FC = () => {
  const dispatch = useDispatch()
  const handleClick = () => {
    dispatch(loadView())
  }

  return (
    <Button
      text={'Submit'}
      color={ComponentColor.Primary}
      onClick={handleClick}
    />
  )
}

export default SubmitQueryButton
