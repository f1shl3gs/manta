// Libraries
import React, {FunctionComponent} from 'react'

// Components
import GetResources from 'src/resources/components/GetResources'
import NotificationEndpointList from 'src/notification_endpoints/components/NotificationEndpointList'
import NotificationEndpointTabHeader from 'src/notification_endpoints/components/NotificationEndpointTabHeader'

// Types
import {ResourceType} from 'src/types/resources'

const NotificationEndpointsIndex: FunctionComponent = () => {
  return (
    <>
      <NotificationEndpointTabHeader />

      <GetResources resources={[ResourceType.NotificationEndpoints]}>
        <NotificationEndpointList />
      </GetResources>
    </>
  )
}

export default NotificationEndpointsIndex
