import React from 'react'
import {Button, ComponentColor} from '@influxdata/clockface'
import {useDispatch} from 'react-redux'
import {poll} from 'src/shared/actions/autoRefresh'

const SubmitQueryButton: React.FC = () => {
  const dispatch = useDispatch()
  const handleClick = () => {
    dispatch(poll())
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
