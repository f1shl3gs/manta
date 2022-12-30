// Libraries
import React, {FunctionComponent, lazy} from 'react'

// Components
import TabsPage from 'src/layout/TabsPage'

const Todo = lazy(() => import('src/Todo'))
const ConfigurationPage = lazy(
  () => import('src/configurations/components/ConfigurationIndex')
)
const ScrapePage = lazy(() => import('src/scrapes/components/ScrapesIndex'))

const tabs = [
  {
    name: 'vertex',
    path: 'vertex',
    element: <Todo />,
  },
  {
    name: 'config',
    path: 'config',
    element: <ConfigurationPage />,
  },
  {
    name: 'scrape',
    path: 'scrape',
    element: <ScrapePage />,
  },
]

const DataPage: FunctionComponent = () => {
  return <TabsPage title={'Data'} tabs={tabs} />
}

export default DataPage
