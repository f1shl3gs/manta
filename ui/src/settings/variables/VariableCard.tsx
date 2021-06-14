// Libraries
import React from 'react'
import moment from 'moment'
import {useHistory} from 'react-router-dom'

// Types
import {Variable} from '../../types/Variable'
import {
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'

// Hooks
import {useOrgID} from '../../shared/useOrg'

interface Props {
  variable: Variable
  onDelete: (v: Variable) => void
  onNameUpdate: (id: string, name: string) => void
  onDescUpdate: (id: string, desc: string) => void
}

const VariableCard: React.FC<Props> = props => {
  const history = useHistory()
  const orgID = useOrgID()
  const {variable, onDelete, onNameUpdate, onDescUpdate} = props
  const context = (id: string, name: string): JSX.Element => {
    return (
      <FlexBox margin={ComponentSize.Small}>
        <FlexBox.Child>
          <Button
            icon={IconFont.Trash}
            text={'Delete'}
            size={ComponentSize.ExtraSmall}
            color={ComponentColor.Danger}
            onClick={() => onDelete(variable)}
          />
        </FlexBox.Child>
      </FlexBox>
    )
  }

  return (
    <ResourceCard
      key={variable.id}
      contextMenu={context(variable.id, variable.name)}
    >
      <ResourceCard.EditableName
        name={variable.name}
        onClick={() => {
          history.push(`/orgs/${orgID}/settings/variables/${variable.id}`)
        }}
        onUpdate={name => onNameUpdate(variable.id, name)}
      />

      <ResourceCard.EditableDescription
        description={variable.desc || ''}
        onUpdate={desc => onDescUpdate(variable.id, desc)}
      />

      <ResourceCard.Meta>
        <span>updated: {moment(variable.updated).fromNow()}</span>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default VariableCard
