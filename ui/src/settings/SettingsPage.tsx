import React, {FunctionComponent, lazy} from 'react'
import TabsPage from 'src/shared/components/TabsPage'

const Members = lazy(() => import('src/settings/Members'))
const Secrets = lazy(() => import('src/settings/Secrets'))

const SettingsPage: FunctionComponent = () => {
  const tabs = [
    {
      name: 'members',
      element: <Members />,
    },
    {
      name: 'secrets',
      element: <Secrets />,
    },
  ]

  return <TabsPage title={'Settings'} tabs={tabs} />
}

export default SettingsPage
