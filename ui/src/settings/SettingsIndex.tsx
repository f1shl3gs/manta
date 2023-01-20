import React, {FunctionComponent, lazy} from 'react'
import TabsPage from 'src/layout/TabsPage'

const Members = lazy(() => import('src/members/components/MembersIndex'))
const Secrets = lazy(() => import('src/secrets/components/SecretsIndex'))

const SettingsPage: FunctionComponent = () => {
  const tabs = [
    {
      name: 'members',
      path: 'members',
      element: <Members />,
    },
    {
      name: 'secrets',
      path: 'secrets',
      element: <Secrets />,
    },
  ]

  return <TabsPage title={'Settings'} tabs={tabs} />
}

export default SettingsPage
