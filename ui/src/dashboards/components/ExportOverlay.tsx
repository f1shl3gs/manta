// Libraries
import React, {FunctionComponent} from 'react'
import {useParams} from 'react-router-dom'
import {useSelector} from 'react-redux'

// Componetns
import ExportForm from 'src/shared/components/ExportForm'
import GetResource from 'src/resources/components/GetResource'

// Types
import {ResourceType} from 'src/types/resources'

// Selectors
import {getDashboardWithCell} from 'src/dashboards/selectors'

const ExportOverlay: FunctionComponent = () => {
  const {dashboardID} = useParams()
  const dashboard = useSelector(getDashboardWithCell(dashboardID))

  return (
    <GetResource resources={[{id: dashboardID, type: ResourceType.Dashboards}]}>
      <ExportForm
        resourceName={'dashboard'}
        name={dashboard.name ?? 'dashboard'}
        content={JSON.stringify(dashboard, null, 2)}
      />
    </GetResource>
  )
}

export default ExportOverlay
