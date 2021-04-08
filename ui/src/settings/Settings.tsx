// Libraries
import React from 'react'

// Components
import {Page, Tabs} from '@influxdata/clockface'

// Hooks
import {useOrgID} from '../shared/useOrg'
import Navigation from '../layout/resourcePage/Navigation'
import {Route, Switch} from 'react-router-dom'
import VariablesPage from './VariablesPage'
import Todo from '../components/Todo'

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
]

const Settings: React.FC = () => {
  const activeColumn = 'secrets'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`
  const pagePrefix = `/orgs/${orgID}/settings`

  return (
    <Page titleTag={`Settings | ${activeColumn}`}>
      <Page.Header fullWidth={true}>
        <Page.Title title={SETTINGS_PAGE_TITLE} />
      </Page.Header>

      <Page.Contents className={pageContentsClassName}>
        <Navigation prefix={pagePrefix} tabs={tabs} />

        <Tabs.TabContents>
          <Switch>
            <Route path={`${pagePrefix}/variables`} component={VariablesPage} />
            <Route path={`${pagePrefix}/variables/:id`} component={Todo} />
            <Route path={`${pagePrefix}/secrets`} component={Todo} />
          </Switch>
        </Tabs.TabContents>
      </Page.Contents>
    </Page>
  )
}

export default Settings
