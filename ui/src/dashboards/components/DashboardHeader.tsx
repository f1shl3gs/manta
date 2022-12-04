// Libraries
import React from 'react'
import {useParams} from 'react-router-dom'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {ComponentSize, FlexBox, Page} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'
import PresentationModeToggle from 'src/shared/components/PresentationModeToggle'
import TimeRangeDropdown from 'src/shared/components/TimeRangeDropdown'
import CreateCellButton from 'src/dashboards/components/CreateCellButton'
import AutoRefreshDropdown from 'src/shared/components/AutoRefreshDropdown'
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'

// Types
import {AppState} from 'src/types/stores'
import {Dashboard} from 'src/types/dashboards'
import {ResourceType} from 'src/types/resources'

// Actions
import {updateDashboard} from 'src/dashboards/actions/thunks'

// Selectors
import {getByID} from 'src/resources/selectors'

const DashboardHeader = () => {
  const dispatch = useDispatch()
  const {dashboardID} = useParams()
  const {name} = useSelector((state: AppState) =>
    getByID<Dashboard>(state, ResourceType.Dashboards, dashboardID)
  )

  const handleRename = (newName: string) => {
    dispatch(updateDashboard(dashboardID, {name: newName}))
  }

  return (
    <Page.Header fullWidth={true}>
      <RenamablePageTitle
        name={name}
        placeholder={'Name this dashboard'}
        maxLength={68}
        onRename={handleRename}
      />

      <FlexBox margin={ComponentSize.Large}>
        <FlexBox margin={ComponentSize.Small}>
          <CreateCellButton />
          <PresentationModeToggle />
        </FlexBox>

        <FlexBox margin={ComponentSize.Small}>
          <TimeRangeDropdown />
          <AutoRefreshDropdown />
          <AutoRefreshButton />
        </FlexBox>
      </FlexBox>

    </Page.Header>
  )
}

export default DashboardHeader
