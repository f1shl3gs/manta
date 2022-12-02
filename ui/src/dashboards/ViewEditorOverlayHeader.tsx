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
import useEscape from 'src/shared/useEscape'

interface Props {
  onSubmit: () => void
}

const ViewEditorOverlayHeader: FunctionComponent<Props> = ({onSubmit}) => {
  const {cell, onRename, setViewProperties} = useCell()
  const onDismiss = useEscape()

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
            onClick={onSubmit}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default ViewEditorOverlayHeader
