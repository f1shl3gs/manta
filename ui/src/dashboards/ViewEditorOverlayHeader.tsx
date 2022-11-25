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
import ViewTypeDropdown from 'src/dashboards/components/ViewTypeDropdown'
import VisOptionsButton from 'src/dashboards/components/VisOptionsButton'

// Hooks
import {useCell} from 'src/dashboards/useCell'

interface Props {
  onDismiss: () => void
}

const ViewEditorOverlayHeader: FunctionComponent<Props> = ({onDismiss}) => {
  const {cell, onRename, updateCell, setViewProperties} = useCell()

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          onRename={onRename}
          name={cell.name}
          placeholder={'Name this cell'}
          maxLength={68}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <ViewTypeDropdown setViewProperties={setViewProperties} />
          <VisOptionsButton />
        </Page.ControlBarLeft>

        <Page.ControlBarRight>
          <SquareButton
            icon={IconFont.Remove_New}
            size={ComponentSize.Small}
            onClick={onDismiss}
          />

          <SquareButton
            icon={IconFont.CheckMark_New}
            size={ComponentSize.Small}
            color={ComponentColor.Success}
            onClick={updateCell}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default ViewEditorOverlayHeader
