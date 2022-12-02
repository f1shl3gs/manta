// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Page} from '@influxdata/clockface'
import PageSpinner from 'src/shared/components/PageSpinner'
import DashboardEmpty from 'src/dashboards/components/DashboardEmpty'

// Hooks
import useFetch from 'src/shared/useFetch'
import {useParams} from 'react-router-dom'
import {Dashboard} from 'src/types/dashboards'
import DashboardHeader from 'src/dashboards/components/DashboardHeader'
import Cells from 'src/dashboards/components/Cells'

const DashboardPage: FunctionComponent = () => {
  const {dashboardId} = useParams()
  const {data: dashboard, loading} = useFetch<Dashboard>(
    `/api/v1/dashboards/${dashboardId}`
  )

  return (
    <PageSpinner loading={loading}>
      <Page titleTag={`Dashboard | ${dashboard?.name}`}>
        <DashboardHeader />

        <Page.Contents scrollable={true}>
          {dashboard?.cells && dashboard.cells.length !== 0 ? (
            <Cells />
          ) : (
            <DashboardEmpty />
          )}
        </Page.Contents>
      </Page>
    </PageSpinner>
  )
}

export default DashboardPage
