// Libraries
import React from 'react'
import {useHistory} from 'react-router-dom'
import moment from 'moment'

// Components
import {
  Button,
  ButtonShape,
  ComponentColor,
  ComponentSize,
  FlexBox,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import CopyButton from '../../shared/components/CopyButton'

// Types
import {Otcl} from 'types/otcl'

// Hooks
import {useOrgID} from '../../shared/useOrg'
import {useOtcls} from './useOtcls'

interface Props {
  otcl: Otcl
  onDelete: (otcl: Otcl) => void
}

const OtclCard: React.FC<Props> = props => {
  const {otcl, onDelete} = props
  const history = useHistory()
  const orgID = useOrgID()
  const {onNameUpdate, onDescUpdate} = useOtcls()

  const context = (id: string): JSX.Element => {
    return (
      <FlexBox margin={ComponentSize.Small}>
        <FlexBox.Child>
          <Button
            icon={IconFont.Trash}
            text="Delete"
            size={ComponentSize.ExtraSmall}
            color={ComponentColor.Danger}
            onClick={() => onDelete(otcl)}
          />
        </FlexBox.Child>

        <FlexBox.Child>
          <CopyButton
            shape={ButtonShape.Default}
            textToCopy={`${window.location.protocol}//${window.location.hostname}/api/v1/otcls/${id}`}
            buttonText={'Copy Otcl config url'}
            color={ComponentColor.Default}
            contentName={'cn'}
            size={ComponentSize.ExtraSmall}
          />
        </FlexBox.Child>
      </FlexBox>
    )
  }

  return (
    <ResourceCard key={otcl.id} contextMenu={context(otcl.id)}>
      <ResourceCard.EditableName
        name={otcl.name}
        onClick={() => {
          history.push(`/orgs/${orgID}/data/otcls/${otcl.id}`)
        }}
        onUpdate={name => onNameUpdate(otcl.id, name)}
      />
      <ResourceCard.EditableDescription
        description={otcl.desc}
        onUpdate={desc => onDescUpdate(otcl.id, desc)}
      />

      <ResourceCard.Meta>
        <span>updated: {moment(otcl.updated).fromNow()}</span>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default OtclCard
