// Libraries
import React, {FunctionComponent} from 'react'
import {useNavigate} from 'react-router-dom'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

const CreateCheckButton: FunctionComponent = () => {
  const navigate = useNavigate()

  const handleClick = () => {
    navigate(`${window.location.pathname}/new`)
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
