// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  ConfirmationButton,
  FlexBox,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'

// Types
import {NotificationEndpoint} from 'src/types/notificationEndpoints'

// Actions
import {
  deleteNotificationEndpoint,
  NotificationEndpointUpdate,
  patchNotificationEndpoint,
} from 'src/notification_endpoints/actions/thunks'

// Utils
import {fromNow} from 'src/shared/utils/duration'
import {useDispatch} from 'react-redux'
import {useNavigate} from 'react-router-dom'

interface OwnProps {
  notificationEndpoint: NotificationEndpoint
}

const NotificationEndpointCard: FunctionComponent<OwnProps> = ({
  notificationEndpoint,
}) => {
  const dispatch = useDispatch()
  const navigate = useNavigate()

  const handleDelete = () => {
    dispatch(deleteNotificationEndpoint(notificationEndpoint.id))
  }
  const handleUpdate = (upd: NotificationEndpointUpdate) => {
    dispatch(patchNotificationEndpoint(notificationEndpoint.id, upd))
  }

  const contextMenu = (): JSX.Element => (
    <FlexBox margin={ComponentSize.ExtraSmall}>
      <ConfirmationButton
        color={ComponentColor.Colorless}
        icon={IconFont.Trash_New}
        shape={ButtonShape.Square}
        size={ComponentSize.ExtraSmall}
        confirmationLabel={'Delete this dashboard'}
        confirmationButtonText={'Confirm'}
        onConfirm={handleDelete}
        testID={'dashboard-card-context--delete'}
      />
    </FlexBox>
  )

  return (
    <ResourceCard key={notificationEndpoint.id} contextMenu={contextMenu()}>
      <ResourceCard.EditableName
        name={notificationEndpoint.name}
        onClick={() =>
          navigate(`${window.location.pathname}/${notificationEndpoint.id}`)
        }
        onUpdate={name => handleUpdate({name})}
      />

      <ResourceCard.EditableDescription
        description={notificationEndpoint.desc}
        onUpdate={desc => handleUpdate({desc})}
      />

      <ResourceCard.Meta>
        <>Modified: {fromNow(notificationEndpoint.updated)}</>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default NotificationEndpointCard
