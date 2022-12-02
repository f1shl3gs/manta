import React, {FunctionComponent, lazy} from 'react'
import TabsPage from 'src/layout/TabsPage'

const Members = lazy(() => import('src/members/MembersIndex'))
const Secrets = lazy(() => import('src/secrets/SecretsIndex'))

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
