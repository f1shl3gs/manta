// Libraries
import React from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'
import ViewTypeDropdown from 'components/timeMachine/ViewTypeDropdown'
import VisOptionsButton from './VisOptionsButton'

// Hooks
import {useCell} from './useCell'
import {useViewProperties} from 'shared/useViewProperties'

const saveButtonClass = 'veo-header--save-cell-button'

const ViewEditorOverlayHeader: React.FC = () => {
  const history = useHistory()
  const {cell, onRename, updateCell} = useCell()
  const {viewProperties} = useViewProperties()
  const onSave = () =>
    updateCell({
      ...cell,
      viewProperties,
    })
  const onCancel = () => history.goBack()

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={cell!.name}
          onRename={onRename}
          placeholder={'Name this Cell'}
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
            icon={IconFont.Remove}
            onClick={onCancel}
            size={ComponentSize.Small}
          />
          <SquareButton
            className={saveButtonClass}
            icon={IconFont.Checkmark}
            color={ComponentColor.Success}
            size={ComponentSize.Small}
            onClick={onSave}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default ViewEditorOverlayHeader
