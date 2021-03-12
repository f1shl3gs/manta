// Libraries
import React, {useCallback} from 'react'
import {useHistory, useLocation, useParams} from 'react-router-dom'

// Components
import {
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'

// Hooks
import {useCheck} from './useCheck'
import {useOrgID} from '../../shared/useOrg'

const saveButtonClass = 'veo-header--save-cell-button'

const CheckOverlayHeader: React.FC = () => {
  const {check, onRename, onSave} = useCheck()
  const history = useHistory()
  const orgID = useOrgID()

  const onCancel = useCallback(() => {
    history.push(`/orgs/${orgID}/alerts/checks`)
  }, [orgID, history])

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={check.name}
          onRename={onRename}
          placeholder={'Name this Check'}
          maxLength={68}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <div>todo</div>
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
            onClick={() => {
              onSave()
                .then(() => {
                  history.goBack()
                })
                .catch(err => {
                  console.warn(err)
                })
            }}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default CheckOverlayHeader
