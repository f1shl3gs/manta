// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'
import ViewTypeDropdown from 'src/timeMachine/ViewTypeDropdown'
import VisOptionsButton from 'src/dashboards/components/VisOptionsButton'

// Hooks

interface Props {
  name: string
  onSubmit: () => void
  onCancel: () => void
  onRename: (name: string) => void
}

const ViewEditorOverlayHeader: FunctionComponent<Props> = ({
  name,
  onSubmit,
  onRename,
  onCancel,
}) => {
  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          onRename={onRename}
          name={name}
          placeholder={'Name this cell'}
          maxLength={68}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <ViewTypeDropdown />
          <VisOptionsButton />
        </Page.ControlBarLeft>

        <Page.ControlBarRight>
          <SquareButton
            icon={IconFont.Remove_New}
            size={ComponentSize.Small}
            onClick={onCancel}
          />

          <SquareButton
            icon={IconFont.CheckMark_New}
            size={ComponentSize.Small}
            color={ComponentColor.Success}
            onClick={onSubmit}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default ViewEditorOverlayHeader
