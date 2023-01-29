// Libraries
import React, {FunctionComponent} from 'react'

// Components
import HTTPOptions from 'src/notification_endpoints/components/HTTPOptions'

// Types
import {NotificationEndpointType} from 'src/types/notificationEndpoints'

interface Props {
  type: NotificationEndpointType
}

const EndpointOptions: FunctionComponent<Props> = ({type}) => {
  switch (type) {
    case 'http':
      return <HTTPOptions />

    default:
      return <>Unknown Type {type}</>
  }
}

export default EndpointOptions
