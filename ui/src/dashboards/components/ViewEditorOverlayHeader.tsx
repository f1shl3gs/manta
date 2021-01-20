import React, {useCallback} from 'react'

import {
  Button,
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'
import VisOptionsButton from './VisOptionsButton'
import {useHistory} from 'react-router-dom'
import {useCell} from './useCell'
import {useViewProperties} from '../../shared/useViewProperties'
import {Cell} from '../../types/Dashboard'
import {useDashboard} from './useDashboard'
import ViewTypeDropdown from '../../components/timeMachine/ViewTypeDropdown'

const saveButtonClass = 'veo-header--save-cell-button'

interface Props {
  name: string
  onNameSet: (name: string) => void
  onSave: () => void
  onCancel: () => void
}

const onClickOutside = () => console.log('onClickOutside')

const ViewEditorOverlayHeader: React.FC = (props) => {
  const history = useHistory()
  const {cell, updateCell} = useCell()
  const {reload} = useDashboard()
  const {viewProperties} = useViewProperties()

  const onNameSet = useCallback(
    (name: string) => {
      updateCell({
        ...cell,
        viewProperties,
        name: name,
      } as Cell)
    },
    [cell, viewProperties]
  )

  const onSave = useCallback(() => {
    updateCell({
      ...cell,
      viewProperties,
    } as Cell).then(() => {
      history.goBack()
      reload()
    })
  }, [cell, viewProperties])

  const onCancel = () => history.goBack()

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={cell!.name}
          onRename={onNameSet}
          placeholder={'Name this Cell'}
          maxLength={68}
          onClickOutside={onClickOutside}
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
