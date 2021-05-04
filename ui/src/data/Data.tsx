// Libraries
import React from 'react'
import {Route, Switch} from 'react-router-dom'

// Components
import {Page, Tabs} from '@influxdata/clockface'
import Navigation from '../layout/resourcePage/Navigation'
import Todo from '../components/Todo'
import Otcls from './otcls/Otcls'
import OtclOverlay from './otcls/OtclOverlay'
import Scrapers from './scrapers/Scrapers'

// Hooks
import {useOrgID} from '../shared/useOrg'

// Constants
const tabs = [
  {
    id: 'otcls',
    text: 'Otcls',
  },
  {
    id: 'scrapers',
    text: 'Scrapers',
  },
]

const Data: React.FC = () => {
  // todo: rename the classname
  const pageContentsClassName = `alerting-index alerting-index__secrets`
  const orgID = useOrgID()
  const pagePrefix = `/orgs/${orgID}/data`

  return (
    <Page titleTag={`Data | ${'abc'}`}>
      <Page.Header fullWidth={false}>
        <Page.Title title={'Data'} />
      </Page.Header>

      <Page.Contents className={pageContentsClassName}>
        <Navigation prefix={pagePrefix} tabs={tabs} />

        <Tabs.TabContents>
          <Switch>
            <Route
              path={`${pagePrefix}/otcls`}
              component={Otcls}
              exact={true}
            />

            <Route path={`${pagePrefix}/otcls/:id`} component={OtclOverlay} />

            <Route
              path={`${pagePrefix}/scrapers`}
              component={Scrapers}
              exact={true}
            />
            <Route
              path={`${pagePrefix}/scrapers/:id`}
              component={Todo}
              exact={true}
            />
          </Switch>
        </Tabs.TabContents>
      </Page.Contents>
    </Page>
  )
}

export default Data
