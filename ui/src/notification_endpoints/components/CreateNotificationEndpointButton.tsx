// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

const CreateNotificationEndpointButton: FunctionComponent = () => {
  return (
    <Button
      text={'Create Notification Endpoint'}
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      testID={'create-notification-endpoint--button'}
      onClick={() => console.log('create')}
    />
  )
}

export default CreateNotificationEndpointButton
