// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Page} from '@influxdata/clockface'
import {Todo} from 'src/Todo'
import PageSpinner from 'src/shared/components/PageSpinner'
import DashboardEmpty from 'src/dashboards/components/DashboardEmpty'

import useFetch from 'src/shared/useFetch'
import {useParams} from 'react-router-dom'
import { Dashboard } from 'src/types/Dashboard'

const DashboardPage: FunctionComponent = () => {
  const {dashboardId} = useParams()
  const {data: dashboard, loading} = useFetch<Dashboard>(`/api/v1/dashboards/${dashboardId}`)

  console.log(dashboard)

  return (
    <>
      <PageSpinner loading={loading}>
        <Page>
          {dashboard?.cells.length === 0 ? <DashboardEmpty /> : <Todo/>}
        </Page>
      </PageSpinner>
    </>
  )
}

export default DashboardPage
