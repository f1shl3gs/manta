import React, {FunctionComponent, lazy} from 'react'
import TabsPage from 'src/shared/components/TabsPage'

const Todo = lazy(() => import('src/Todo'))
const ConfigurationPage = lazy(
  () => import('src/data/configuration/ConfigurationPage')
)

const DataPage: FunctionComponent = () => {
  const tabs = [
    {
      name: 'vertex',
      element: <Todo />,
    },
    {
      name: 'config',
      element: <ConfigurationPage />,
    },
  ]

  return (
    <>
      <TabsPage title={'Data'} tabs={tabs} />
    </>
  )
}

export default DataPage
