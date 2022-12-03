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
import VisOptionsButton from 'src/timeMachine/VisOptionsButton'

// Hooks
import useEscape from 'src/shared/useEscape'
import {useTimeMachine} from 'src/timeMachine/useTimeMachine'

interface Props {
  name: string
  onSubmit: () => void
  onRename: (name: string) => void
}

const ViewEditorOverlayHeader: FunctionComponent<Props> = ({
  name,
  onRename,
  onSubmit,
}) => {
  const {setViewProperties} = useTimeMachine()
  const onDismiss = useEscape()

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
