// Libraries
import React from 'react'

// Components
import {
  AlignItems,
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  FlexDirection,
  IconFont,
  JustifyContent,
  ResourceCard,
} from '@influxdata/clockface'

// Hooks
import {useHistory} from 'react-router-dom'
import {useOrgID} from '../../shared/useOrg'

// Utils
import {relativeTimestampFormatter} from '../../utils/relativeTimestampFormatter'

// Types
import {NotificationEndpoint} from '../../client'

interface Props {
  endpoint: NotificationEndpoint
}

const NotificationEndpointCard: React.FC<Props> = props => {
  const orgID = useOrgID()
  const history = useHistory()

  const {
    endpoint: {id, name, desc, updated},
  } = props

  const onNameUpdate = (v: string) => {
    console.log('new name')
  }

  const onDescUpdate = (v: string) => {
    console.log('new desc')
  }

  const contextMenu = () => (
    <Button
      icon={IconFont.Trash}
      text={'Delete'}
      color={ComponentColor.Danger}
      size={ComponentSize.ExtraSmall}
      onClick={() => {
        console.log('delete')
      }}
    />
  )

  return (
    <ResourceCard
      key={`notification-endpoint-id--${id}`}
      direction={FlexDirection.Row}
      alignItems={AlignItems.Center}
      margin={ComponentSize.Large}
      contextMenu={contextMenu()}
    >
      <FlexBox
        direction={FlexDirection.Column}
        justifyContent={JustifyContent.Center}
        margin={ComponentSize.Medium}
        alignItems={AlignItems.FlexStart}
      >
        todo
      </FlexBox>

      <FlexBox
        direction={FlexDirection.Column}
        margin={ComponentSize.Small}
        alignItems={AlignItems.FlexStart}
      >
        <ResourceCard.EditableName
          name={name}
          noNameString={'Name this Endpoint'}
          onUpdate={onNameUpdate}
          onClick={() => {
            history.push(`/orgs/${orgID}/alerts/endpoints/${id}`)
          }}
        />

        <ResourceCard.EditableDescription
          description={desc}
          placeholder={`Describe ${name}`}
          onUpdate={onDescUpdate}
        />

        <ResourceCard.Meta>
          <>Last completed at blah blah</>
          <>{relativeTimestampFormatter(updated, 'Last updated ')}</>
        </ResourceCard.Meta>
      </FlexBox>
    </ResourceCard>
  )
}

export default NotificationEndpointCard