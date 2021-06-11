// Libraries
import React, {useCallback} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
  SelectGroup,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'components/RenamablePageTitle'

// Hooks
import {useCheck} from './useCheck'
import {useOrgID} from '../../shared/useOrg'
import {
  defaultErrorNotification,
  useNotification,
} from '../../shared/notification/useNotification'

const saveButtonClass = 'veo-header--save-cell-button'

const CheckOverlayHeader: React.FC = () => {
  const {id, name, onRename, onSave, tab, setTab} = useCheck()
  const history = useHistory()
  const orgID = useOrgID()
  const {notify} = useNotification()

  const onCancel = useCallback(() => {
    history.push(`/orgs/${orgID}/alerts/checks`)
  }, [orgID, history])

  // todo: it's a bad name
  const onSubmit = useCallback(() => {
    onSave()
      .then(() => {
        console.log('success')
        // history.push(`/orgs/${orgID}/alerts/checks`)
      })
      .catch(err => {
        const message =
          id === 'new'
            ? `Create new checks failed, err: ${err.message}`
            : `Update Check ${name} failed, err: ${err}`

        notify({
          ...defaultErrorNotification,
          message,
        })
      })
  }, [onSave, history, orgID, id, name, notify])

  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          name={name}
          onRename={onRename}
          placeholder={'Name this Check'}
          maxLength={68}
        />
      </Page.Header>

      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <SelectGroup
            shape={ButtonShape.StretchToFit}
            style={{width: '300px'}}
          >
            <SelectGroup.Option
              id={'query'}
              value={'query'}
              active={tab === 'query'}
              onClick={setTab}
            >
              1. Define Query
            </SelectGroup.Option>
            <SelectGroup.Option
              id={'check'}
              value={'check'}
              active={tab === 'check'}
              onClick={setTab}
            >
              2. Configure Check
            </SelectGroup.Option>
          </SelectGroup>
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
            onClick={onSubmit}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  )
}

export default CheckOverlayHeader
