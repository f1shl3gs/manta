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

const saveButtonClass = 'veo-header--save-cell-button'

const CheckOverlayHeader: React.FC = () => {
  const {name, onRename, onSave} = useCheck()
  const history = useHistory()
  const orgID = useOrgID()

  const onCancel = useCallback(() => {
    history.push(`/orgs/${orgID}/alerts/checks`)
  }, [orgID, history])

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
              id={'A'}
              value={'a'}
              active
              onClick={v => console.log(v)}
            >
              Define Query
            </SelectGroup.Option>
            <SelectGroup.Option
              id={'B'}
              value={'b'}
              active={false}
              onClick={v => console.log(v)}
            >
              Configure Check
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
