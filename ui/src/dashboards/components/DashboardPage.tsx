// Libraries
import React, {FunctionComponent, lazy, useEffect} from 'react'
import {Route, Routes, useParams} from 'react-router-dom'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {Page} from '@influxdata/clockface'
import DashboardEmpty from 'src/dashboards/components/DashboardEmpty'
import DashboardHeader from 'src/dashboards/components/DashboardHeader'
import Cells from 'src/dashboards/components/Cells'
import GetResource from 'src/resources/components/GetResource'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'
import {Dashboard} from 'src/types/dashboards'

// Actions
import {poll, setAutoRefreshInterval} from 'src/shared/actions/autoRefresh'
import {setTimeRange} from 'src/shared/actions/timeRange'

// Selectors
import {getByID} from 'src/resources/selectors'
import {getCells} from 'src/cells/selectors'

// Constants
import {pastHourTimeRange} from 'src/constants/timeRange'

// Lazy Loads
const NewVEO = lazy(() => import('src/dashboards/components/NewVEO'))
const EditVEO = lazy(() => import('src/dashboards/components/EditVEO'))

interface Props {
  id: string
}

const DashboardIndex: FunctionComponent<Props> = ({id}) => {
  const dispatch = useDispatch()
  const {name, cells} = useSelector((state: AppState) => {
    const dashbaord = getByID<Dashboard>(state, ResourceType.Dashboards, id)
    const cells = getCells(state, dashbaord.id)

    return {
      name: dashbaord.name ?? '',
      cells: cells ?? [],
    }
  })

  useEffect(() => {
    dispatch(setTimeRange(pastHourTimeRange))
    dispatch(setAutoRefreshInterval(15))
    dispatch(poll())
  }, [dispatch])

  return (
    <>
      <Page titleTag={`Dashboard | ${name}`}>
        <DashboardHeader />

        <Page.Contents scrollable={true}>
          {cells.length !== 0 ? <Cells /> : <DashboardEmpty />}
        </Page.Contents>
      </Page>
    </>
  )
}

export default () => {
  const {dashboardID} = useParams()

  return (
    <GetResource resources={[{type: ResourceType.Dashboards, id: dashboardID}]}>
      <Routes>
        <Route index element={<DashboardIndex id={dashboardID} />} />
        <Route path="cells/new" element={<NewVEO />} />
        <Route path="cells/:cellID/edit" element={<EditVEO />} />
      </Routes>
    </GetResource>
  )
}
