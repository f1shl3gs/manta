// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  ConfirmationButton,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'

// Types
import {User} from 'src/types/user'

// Utils
import {fromNow} from 'src/shared/utils/duration'

interface Props {
  member: User
}

const MemberCard: FunctionComponent<Props> = ({member}) => {
  const lastSeen = member.lastSeen ? fromNow(member.lastSeen) : 'Never'
  const handleDelete = () => {
    console.log('delete')
  }

  const contextMenu = (): JSX.Element => (
    <ConfirmationButton
      color={ComponentColor.Colorless}
      icon={IconFont.Trash_New}
      shape={ButtonShape.Square}
      size={ComponentSize.ExtraSmall}
      confirmationLabel={'Delete this user'}
      confirmationButtonText={'Confirm'}
      onConfirm={handleDelete}
      testID={'member-card-context--delete'}
    />
  )

  return (
    <ResourceCard key={member.id} contextMenu={contextMenu()}>
      <ResourceCard.EditableName
        onUpdate={() => console.log('name')}
        name={member.name}
        noNameString={''}
        buttonTestID="editable-name"
        inputTestID="input-field"
      />

      <ResourceCard.Meta>
        <>Last updated: {fromNow(member.updated)} </>
        <>Last Seen: {lastSeen}</>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default MemberCard
