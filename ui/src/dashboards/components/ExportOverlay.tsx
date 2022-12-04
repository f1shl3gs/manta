import React, {FunctionComponent} from 'react'
import ExportForm from 'src/shared/components/ExportForm'
import {useParams} from 'react-router-dom'
import useFetch from 'src/shared/useFetch'
import {Dashboard} from 'src/types/dashboards'
import PageSpinner from 'src/shared/components/PageSpinner'

const ExportOverlay: FunctionComponent = () => {
  const {dashboardID} = useParams()
  const {data, loading} = useFetch<Dashboard>(
    `/api/v1/dashboards/${dashboardID}`
  )

  return (
    <PageSpinner loading={loading}>
      <ExportForm
        resourceName={'dashboard'}
        name={data?.name ?? 'dashboard'}
        content={JSON.stringify(data, null, 2)}
      />
    </PageSpinner>
  )
}

export default ExportOverlay