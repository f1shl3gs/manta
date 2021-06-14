// Libraries
import React from 'react'
import {Route, Switch} from 'react-router-dom'

// Components
import {Page, Tabs} from '@influxdata/clockface'
import Navigation from '../layout/resourcePage/Navigation'
import Variables from './variables/Variables'
import Todo from '../components/Todo'

// Hooks
import {useOrgID} from '../shared/useOrg'

const SETTINGS_PAGE_TITLE = 'Settings'

const tabs = [
  {
    id: 'variables',
    text: 'Variables',
  },
  {
    id: 'secrets',
    text: 'Secrets',
  },
  {
    id: 'templates',
    text: 'Templates',
  },
]

const Settings: React.FC = () => {
  const activeColumn = 'secrets'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`
  const pagePrefix = `/orgs/${orgID}/settings`

  return (
    <Page titleTag={`Settings | ${activeColumn}`}>
      <Page.Header fullWidth={false}>
        <Page.Title title={SETTINGS_PAGE_TITLE} />
      </Page.Header>

      <Page.Contents className={pageContentsClassName} fullWidth={false}>
        <Navigation prefix={pagePrefix} tabs={tabs} />

        <Tabs.TabContents>
          <Switch>
            <Route path={`${pagePrefix}/variables/:id`} component={Todo} />
            <Route path={`${pagePrefix}/variables`} component={Variables} />
            <Route path={`${pagePrefix}/secrets/:id`} component={Todo} />
            <Route path={`${pagePrefix}/secrets`} component={Todo} />
          </Switch>
        </Tabs.TabContents>
      </Page.Contents>
    </Page>
  )
}

export default Settings
