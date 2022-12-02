// Libraries
import React from 'react'

// Components
import {ComponentSize, FlexBox, Page} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'
import {useDashboard} from 'src/dashboards/useDashboard'
import PresentationModeToggle from 'src/dashboards/components/PresentationModeToggle'
import TimeRangeDropdown from 'src/dashboards/components/TimeRangeDropdown'
import CreateCellButton from 'src/dashboards/components/CreateCellButton'
import AutoRefreshDropdown from 'src/shared/components/AutoRefreshDropdown'
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'

const DashboardHeader = () => {
  const {name, onRename} = useDashboard()

  return (
    <Page.Header fullWidth={true}>
      <RenamablePageTitle
        name={name}
        placeholder={'Name this dashboard'}
        maxLength={68}
        onRename={onRename}
      />

      <FlexBox margin={ComponentSize.Small}>
        <CreateCellButton />
        <PresentationModeToggle />
        <TimeRangeDropdown />
      </FlexBox>

      <FlexBox margin={ComponentSize.Small}>
        <AutoRefreshDropdown />
        <AutoRefreshButton />
      </FlexBox>
    </Page.Header>
  )
}

export default DashboardHeader
