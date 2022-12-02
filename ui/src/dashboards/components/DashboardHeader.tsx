// Libraries
import React from 'react'

// Components
import {ComponentSize, FlexBox, Page} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'
import PresentationModeToggle from 'src/shared/components/PresentationModeToggle'
import TimeRangeDropdown from 'src/shared/components/TimeRangeDropdown'
import CreateCellButton from 'src/dashboards/components/CreateCellButton'
import AutoRefreshDropdown from 'src/shared/components/AutoRefreshDropdown'
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'

// Hooks
import {useDispatch, useSelector} from 'react-redux'

// Actions
import {updateDashboard} from 'src/dashboards/actions/thunks'

// Selectors
import {getDashboard} from 'src/dashboards/selectors'

const DashboardHeader = () => {
  const {id, name} = useSelector(getDashboard)
  const dispatch = useDispatch()

  const handleRename = (newName: string) => {
    dispatch(updateDashboard(id, {name: newName}))
  }

  return (
    <Page.Header fullWidth={true}>
      <RenamablePageTitle
        name={name}
        placeholder={'Name this dashboard'}
        maxLength={68}
        onRename={handleRename}
      />

      <FlexBox margin={ComponentSize.Small}>
        <CreateCellButton />
        <PresentationModeToggle />
      </FlexBox>

      <FlexBox margin={ComponentSize.Small}>
        <TimeRangeDropdown />
        <AutoRefreshDropdown />
        <AutoRefreshButton />
      </FlexBox>
    </Page.Header>
  )
}

export default DashboardHeader
