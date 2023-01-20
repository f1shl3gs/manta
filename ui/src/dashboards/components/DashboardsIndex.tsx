// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Page, PageContents, PageHeader, PageTitle} from '@influxdata/clockface'
import GetResources from 'src/resources/components/GetResources'
import DashboardList from 'src/dashboards/components/DashboardList'
import DashboardsPageHeader from 'src/dashboards/components/DashboardsPageHeader'

// Types
import {ResourceType} from 'src/types/resources'

export const DashboardsIndex: FunctionComponent = () => {
  return (
    <>
      <Page titleTag="Dashboards">
        <PageHeader fullWidth={false}>
          <PageTitle title="Dashboards" />
        </PageHeader>

        <DashboardsPageHeader />

        <PageContents>
          <GetResources resources={[ResourceType.Dashboards]}>
            <DashboardList />
          </GetResources>
        </PageContents>
      </Page>
    </>
  )
}

export default DashboardsIndex
