// Libraries
import React from 'react'
import {Route, Switch} from 'react-router-dom'

// Components
import {Page, Tabs} from '@influxdata/clockface'
import AlertsNavigation from './AlertsNavigation'
import Checks from './checks/Checks'
import NotificationEndpoints from './notificationEndpoints/NotificationEndpoints'

// Hooks
import {useOrgID} from 'shared/useOrg'
import CheckOverlay from './checks/CheckOverlay'
import Todo from '../components/Todo'

const ALERTS_PAGE_TITLE = 'Alerts'

const tabs = [
  {
    id: 'checks',
    text: 'Checks',
  },
  {
    id: 'endpoints',
    text: 'Notification Endpoints',
  },
]

const AlertsPage: React.FC = () => {
  // todo: handle activeColumn
  const activeColumn = 'checks'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`
  const pagePrefix = `/orgs/${orgID}/alerts`

  return (
    <Page titleTag={'Alerts | Checks'}>
      <Page.Header fullWidth={true}>
        <Page.Title title={ALERTS_PAGE_TITLE} />
      </Page.Header>

      <Page.Contents className={pageContentsClassName} fullWidth={true}>
        <AlertsNavigation prefix={`${pagePrefix}`} tabs={tabs} />
        <Tabs.TabContents>
          <Switch>
            <Route
              path={`${pagePrefix}/checks/:id`}
              component={CheckOverlay}
              exact={false}
            />

            <Route path={`${pagePrefix}/checks`} component={Checks} />
            <Route
              path={`${pagePrefix}/endpoints`}
              component={NotificationEndpoints}
            />
          </Switch>
        </Tabs.TabContents>
      </Page.Contents>
    </Page>
  )
}

export default AlertsPage
