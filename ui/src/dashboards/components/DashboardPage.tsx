// Libraries
import React, {FunctionComponent, useEffect} from 'react'

// Components
import {Page} from '@influxdata/clockface'
import DashboardEmpty from 'src/dashboards/components/DashboardEmpty'

// Hooks
import {useParams} from 'react-router-dom'
import DashboardHeader from 'src/dashboards/components/DashboardHeader'
import Cells from 'src/dashboards/components/Cells'
import GetResource from 'src/resources/components/GetResource'
import {useDispatch, useSelector} from 'react-redux'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'
import {Dashboard} from 'src/types/dashboards'

// Selectors
import {getByID} from 'src/resources/selectors'
import {poll, setAutoRefreshInterval} from 'src/shared/actions/autoRefresh'
import {pastHourTimeRange} from 'src/constants/timeRange'
import {setTimeRange} from 'src/shared/actions/timeRange'

interface Props {
  id: string
}

const DashboardIndex: FunctionComponent<Props> = ({id}) => {
  const dispatch = useDispatch()
  const {name, cells} = useSelector((state: AppState) => {
    const dashbaord = getByID<Dashboard>(state, ResourceType.Dashboards, id)

    return {
      name: dashbaord.name ?? '',
      cells: dashbaord.cells ?? [],
    }
  })

  useEffect(() => {
    dispatch(setTimeRange(pastHourTimeRange))
    dispatch(setAutoRefreshInterval(15))
    dispatch(poll())
  }, [dispatch])

  return (
    <Page titleTag={`Dashboard | ${name}`}>
      <DashboardHeader />

      <Page.Contents scrollable={true}>
        {cells.length !== 0 ? <Cells /> : <DashboardEmpty />}
      </Page.Contents>
    </Page>
  )
}

export default () => {
  const {dashboardID} = useParams()

  return (
    <GetResource resources={[{type: ResourceType.Dashboards, id: dashboardID}]}>
      <DashboardIndex id={dashboardID} />
    </GetResource>
  )
}
