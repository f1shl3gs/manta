// Libraries
import React from 'react'

// Components
import {ComponentSize, FlexBox, Page} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'
import {useDashboard} from 'src/dashboards/useDashboard'
import PresentationModeToggle from 'src/dashboards/components/PresentationModeToggle'
import TimeRangeDropdown from 'src/dashboards/components/TimeRangeDropdown'
import CreateCellButton from 'src/dashboards/components/CreateCellButton'

const DashboardHeader = () => {
  const {name} = useDashboard()

  return (
    <Page.Header fullWidth={true}>
      <RenamablePageTitle
        name={name}
        placeholder={'Name this dashboard'}
        maxLength={90}
        onRename={n => console.log(n)}
      />

      <FlexBox margin={ComponentSize.ExtraSmall}>
        <CreateCellButton />
        <PresentationModeToggle />
        <TimeRangeDropdown />
      </FlexBox>
    </Page.Header>
  )
}

export default DashboardHeader
