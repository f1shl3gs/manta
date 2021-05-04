// Libraries
import React from 'react'

// Components
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from '../../components/RenamablePageTitle'

// Hooks
import {useOtcl} from './useOtcl'

const saveButtonClass = 'veo-header--save-cell-button'

interface Props {
  onDismiss: () => void
}

const OtclOverlayHeader: React.FC<Props> = props => {
  const {otcl, onRename, onSave} = useOtcl()
  const {name, content} = otcl
  const {onDismiss} = props

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={name}
          onRename={onRename}
          placeholder={'Name this Otcl'}
          maxLength={68}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarRight>
          <SquareButton
            icon={IconFont.Remove}
            size={ComponentSize.Small}
            onClick={onDismiss}
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

export default OtclOverlayHeader
