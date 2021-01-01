import React from 'react'

import {
  Button,
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'
import VisOptionsButton from './VisOptionsButton';

const saveButtonClass = 'veo-header--save-cell-button'

interface Props {
  name: string
  onNameSet: (name: string) => void
  onSave: () => void
  onCancel: () => void
}

const onClickOutside = () => console.log('onClickOutside')

const ViewEditorOverlayHeader: React.FC<Props> = (props) => {
  const {name, onNameSet, onSave, onCancel} = props

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={name}
          onRename={onNameSet}
          placeholder={'Name this Cell'}
          maxLength={68}
          onClickOutside={onClickOutside}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <Button text={'ViewTypeDropdown'} />
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
