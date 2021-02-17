// Libraries
import React from 'react'

// Components
import {Orientation, Page, Tabs} from '@influxdata/clockface'
import AlertsNavigation from './AlertsNavigation'
import {useOrgID} from '../shared/useOrg'

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

const AlertsPage: React.FC = (props) => {
  const {children} = props
  const activeColumn = 'checks'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`

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
        <AlertsNavigation prefix={`/${orgID}/alerts`} tabs={tabs} />
        <Tabs.Container orientation={Orientation.Horizontal}>
          {children}
        </Tabs.Container>
      </Page.Contents>
    </Page>
  )
}

export default AlertsPage
