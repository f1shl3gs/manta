// Libraries
import React from 'react'

// Components
import {Page} from '@influxdata/clockface'

// Hooks
import {useOrgID} from '../shared/useOrg'

const SETTINGS_PAGE_TITLE = 'Settings'

const tabs = [
  {
    id: 'secrets',
    text: 'Secrets',
  },
]

const SettingsPage: React.FC = () => {
  const activeColumn = 'secrets'
  const orgID = useOrgID()
  const pageContentsClassName = `alerting-index alerting-index__${activeColumn}`
  const pagePrefix = `/orgs/${orgID}/secrets`

  return (
    <Page titleTag={`Settings | ${activeColumn}`}>
      <Page.Header fullWidth={true}>
        <Page.Title title={SETTINGS_PAGE_TITLE} />
      </Page.Header>

      <Page.Contents className={pageContentsClassName}></Page.Contents>
    </Page>
  )
}

export default SettingsPage
