// Libraries
import React, {FunctionComponent, lazy} from 'react'
import {Route, Routes} from 'react-router-dom'

// Components
import {Page, PageContents, PageHeader, PageTitle} from '@influxdata/clockface'
import GetResources from 'src/resources/components/GetResources'
import DashboardList from 'src/dashboards/components/DashboardList'
import DashboardsPageHeader from 'src/dashboards/components/DashboardsPageHeader'

// Types
import {ResourceType} from 'src/types/resources'

// Lazy load overlay
const DashboardImportOverlay = lazy(
  () => import('src/dashboards/components/DashboardImportOverlay')
)

export const DashboardsIndex: FunctionComponent = () => {
  return (
    <>
      <Routes>
        <Route path="import" element={<DashboardImportOverlay />} />
      </Routes>

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
