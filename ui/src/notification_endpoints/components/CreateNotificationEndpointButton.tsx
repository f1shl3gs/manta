// Libraries
import React, {FunctionComponent} from 'react'
import {useNavigate} from 'react-router-dom'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

const CreateNotificationEndpointButton: FunctionComponent = () => {
  const navigate = useNavigate()
  const handleClick = () => {
    navigate(`${window.location.pathname}/new`)
  }

  return (
    <Button
      text={'Create Notification Endpoint'}
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      testID={'create-notification-endpoint--button'}
      onClick={handleClick}
    />
  )
}

export default CreateNotificationEndpointButton
