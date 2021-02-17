// Libraries
import React from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {Button, ComponentColor, IconFont, Page} from '@influxdata/clockface'
import DashboardCards from './DashboardCards'
import SearchWidget from '../shared/components/SearchWidget'

// Hooks
import {DashboardsProvider} from './useDashboards'
import {useFetch} from 'use-http'
import {useOrgID} from 'shared/useOrg'

// Types
import {Dashboard} from 'types/Dashboard'

const useCreateDash = () => {
  const orgID = useOrgID()

  const {post} = useFetch<Dashboard>(`/api/v1/dashboards?orgID=${orgID}`, {
    body: {
      orgID,
    },
  })

  return post
}

const DashboardsIndex: React.FC = () => {
  const create = useCreateDash()
  const history = useHistory()
  const orgID = useOrgID()

  return (
    <Page titleTag={'Dashboards'}>
      <Page.Header fullWidth={false}>
        <Page.Title title={'Dashboards'} />
        {/* rateLimitAlert? */}
      </Page.Header>

      <Page.ControlBar fullWidth={false}>
        <Page.ControlBarLeft>
          <SearchWidget
            search={'v'}
            placeholder={'Filter dashboards...'}
            onSearch={(v) => console.log('v', v)}
          />
        </Page.ControlBarLeft>
        <Page.ControlBarRight>
          <Button
            text={'Add'}
            icon={IconFont.Plus}
            color={ComponentColor.Primary}
            onClick={() => {
              create()
                .then((resp) => {
                  const path = `/orgs/${orgID}/dashboards/${resp.id}`
                  console.log('path', path)
                  history.push(path)
                })
                .catch((err) => {
                  console.log('create dashboard failed', err)
                })
            }}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>

      <Page.Contents
        className="dashboards-index__page-contents"
        fullWidth={false}
        scrollable={true}
      >
        <DashboardsProvider>
          <DashboardCards />
        </DashboardsProvider>
      </Page.Contents>
    </Page>
  )
}

export default DashboardsIndex
