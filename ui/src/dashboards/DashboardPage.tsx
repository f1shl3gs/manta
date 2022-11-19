// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Page} from '@influxdata/clockface'
import PageSpinner from 'src/shared/components/PageSpinner'
import DashboardEmpty from 'src/dashboards/components/DashboardEmpty'

import useFetch from 'src/shared/useFetch'
import {useParams} from 'react-router-dom'
import {Dashboard} from 'src/types/Dashboard'
import {DashboardProvider} from './useDashboard'
import DashboardHeader from './components/DashboardHeader'
import {TimeRangeProvider} from '../shared/useTimeRange'
import Cells from './components/Cells'

const DashboardPage: FunctionComponent = () => {
  const {dashboardId} = useParams()
  const {data: dashboard, loading} = useFetch<Dashboard>(
    `/api/v1/dashboards/${dashboardId}`
  )

  return (
    <PageSpinner loading={loading}>
      <DashboardProvider dashboard={dashboard!}>
        <Page titleTag={`Dashboard | ${dashboard?.name}`}>
          <TimeRangeProvider>
            <DashboardHeader />

            <Page.Contents>
              {dashboard?.cells && dashboard.cells.length !== 0 ? (
                <Cells />
              ) : (
                <DashboardEmpty />
              )}
            </Page.Contents>
          </TimeRangeProvider>
        </Page>
      </DashboardProvider>
    </PageSpinner>
  )
}

export default DashboardPage
