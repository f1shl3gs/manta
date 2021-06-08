// Libraries
import React from 'react'

// Components
import {
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  IconFont,
  Page,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'
import PresentationModeToggle from './PresentationModeToggle'
import AutoRefreshDropdown from 'components/AutoRefreshDropdown/AutoRefreshDropdown'
import TimeRangeDropdown from './TimeRangeDropdown'

// Hooks
import {useDashboard} from './useDashboard'

// Constants
import {AutoRefreshDropdownOptions} from 'constants/autoRefresh'

const DashboardHeader = () => {
  const {
    name,
    onRename,
    addCell,
    showVariablesControls,
    toggleShowVariablesControls,
  } = useDashboard()

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          placeholder={'Name this dashboard'}
          name={name}
          maxLength={90}
          onRename={onRename}
        />

        <FlexBox margin={ComponentSize.ExtraSmall}>
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
          <Button
            icon={IconFont.Cube}
            text={'Variables'}
            testID={'variables--button'}
            onClick={toggleShowVariablesControls}
            color={
              showVariablesControls
                ? ComponentColor.Secondary
                : ComponentColor.Default
            }
          />
          <PresentationModeToggle />

          <TimeRangeDropdown />
          <AutoRefreshDropdown options={AutoRefreshDropdownOptions} />
        </FlexBox>
      </Page.Header>
    </>
  )
}

export default DashboardHeader
