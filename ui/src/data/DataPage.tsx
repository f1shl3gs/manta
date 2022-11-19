// Libraries
import React, {FunctionComponent, lazy} from 'react'

// Components
import TabsPage from 'src/shared/components/TabsPage'

const Todo = lazy(() => import('src/Todo'))
const ConfigurationPage = lazy(
  () => import('src/data/configuration/ConfigurationPage')
)
const ScrapePage = lazy(() => import('src/data/scrape/ScrapePage'))

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
    {
      name: 'scrape',
      element: <ScrapePage />,
    },
  ]

  return (
    <>
      <TabsPage title={'Data'} tabs={tabs} />
    </>
  )
}

export default DataPage
