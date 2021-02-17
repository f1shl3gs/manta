// Libraries
import React from 'react'
import {Route, Switch} from 'react-router-dom'

// Components
import {Page, Tabs} from '@influxdata/clockface'
import AlertsNavigation from './AlertsNavigation'
import ChecksIndex from './ChecksIndex'

// Hooks
import {useOrgID} from 'shared/useOrg'
import {ChecksProvider} from './useChecks'

const ALERTS_PAGE_TITLE = 'Alerts'

const tabs = [
  {
    id: 'checks',
    text: 'Checks',
  },
  {
    id: 'notificationEndpoints',
    text: 'Notification Endpoints',
  },
]

const dummy: React.FC = () => {
  return <div>Dummy</div>
}

const AlertsPage: React.FC = (props) => {
  const activeColumn = 'checks'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`
  const pagePrefix = `/orgs/${orgID}/alerts`

  return (
    <Page titleTag={'Alerts | Checks'}>
      <Page.Header fullWidth={true}>
        <Page.Title title={ALERTS_PAGE_TITLE} />
      </Page.Header>

      <Page.Contents
        fullWidth={true}
        scrollable={false}
        className={pageContentsClassName}
      >
        <AlertsNavigation prefix={`${pagePrefix}`} tabs={tabs} />
        <Tabs.TabContents>
          <Switch>
            <ChecksProvider>
              <Route path={`${pagePrefix}/checks`} component={ChecksIndex} />
            </ChecksProvider>
            <Route
              path={`${pagePrefix}/notificationEndpoints`}
              component={dummy}
            />
          </Switch>
        </Tabs.TabContents>
      </Page.Contents>
    </Page>
  )
}

export default AlertsPage
