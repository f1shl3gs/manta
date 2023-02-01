// Libraries
import React from 'react'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {
  ComponentSize,
  Gradients,
  IconFont,
  Notification,
} from '@influxdata/clockface'

// Types
import {NotificationStyle} from 'src/types/notification'
import {AppState} from 'src/types/stores'

// Actions
import {dismissNotification} from 'src/shared/actions/notifications'

// Utils
import {get} from 'src/shared/utils/get'

const matchGradientToColor = (style: NotificationStyle): Gradients => {
  const converter = {
    [NotificationStyle.Primary]: Gradients.Info,
    [NotificationStyle.Warning]: Gradients.WarningLight,
    [NotificationStyle.Success]: Gradients.HotelBreakfast,
    [NotificationStyle.Error]: Gradients.DangerDark,
    [NotificationStyle.Info]: Gradients.DefaultLight,
  }

  return get(converter, style, Gradients.DefaultLight)
}

const Notifications = () => {
  const notifications = useSelector((state: AppState) => state.notifications)
  const dispatch = useDispatch()

  return (
    <>
      {notifications.map(item => {
        const {id, message, style, duration} = item
        const gradient = matchGradientToColor(style)
        const dismiss = () => {
          dispatch(dismissNotification(id))
        }

        return (
          <Notification
            key={id}
            id={id}
            icon={IconFont.Remove_New}
            duration={duration}
            size={ComponentSize.Small}
            gradient={gradient}
            onDismiss={dismiss}
            onTimeout={dismiss}
            testID={`notification-${style}`}
          >
            <span className={'notification--message'}>{message}</span>
          </Notification>
        )
      })}
    </>
  )
}

export default Notifications
