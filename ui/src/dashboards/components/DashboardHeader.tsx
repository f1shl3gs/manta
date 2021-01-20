// Libraries
import React from 'react'

// Components
import {Button, ComponentColor, IconFont, Page} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'
import PresentationModeToggle from './PresentationModeToggle'
import AutoRefreshDropdown from 'components/AutoRefreshDropdown/AutoRefreshDropdown'
import TimeRangeDropdown from './TimeRangeDropdown'

// Hooks
import {useDashboard} from './useDashboard'

// Constants
import {AutoRefreshDropdownOptions} from 'constants/autoRefresh'

const DashboardHeader = () => {
  const {name, onRename, addCell} = useDashboard()

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          placeholder={'Name this dashboard'}
          name={name}
          maxLength={90}
          onRename={onRename}
        />
      </Page.Header>
      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <Button
            text={'Add Cell'}
            color={ComponentColor.Primary}
            icon={IconFont.AddCell}
            onClick={addCell}
          />
          <Button
            icon={IconFont.TextBlock}
            text="Add Note"
            onClick={() => console.log('add note')}
            testID="add-note--button"
          />
          <PresentationModeToggle />
        </Page.ControlBarLeft>

        <Page.ControlBarRight>
          <TimeRangeDropdown />
          <AutoRefreshDropdown options={AutoRefreshDropdownOptions} />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default DashboardHeader
