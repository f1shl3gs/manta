// Libraries
import React, {FunctionComponent} from 'react'

// Components
import TabsPage from 'src/layout/TabsPage'
import ChecksIndex from 'src/checks/components/ChecksIndex'
import NotificationEndpointsIndex from 'src/notification_endpoints/components/NotificationEndpointsIndex'

const tabs = [
  {
    name: 'Checks',
    path: 'checks',
    element: <ChecksIndex />,
  },
  {
    name: 'Notification Endpoints',
    path: 'notificationEndpoints',
    element: <NotificationEndpointsIndex />,
  },
]

const AlertsIndex: FunctionComponent = () => {
  return <TabsPage title={'Alerts'} tabs={tabs} />
}

export default AlertsIndex
