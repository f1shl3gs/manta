// Libraries
import React from 'react'
import {get} from 'lodash'

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
import {useSelector} from 'react-redux'

const matchGradientToColor = (style: NotificationStyle): Gradients => {
  const converter = {
    [NotificationStyle.Primary]: Gradients.Info,
    [NotificationStyle.Warning]: Gradients.WarningLight,
    [NotificationStyle.Success]: Gradients.HotelBreakfast,
    [NotificationStyle.Error]: Gradients.DangerDark,
    [NotificationStyle.Info]: Gradients.DefaultLight,
  }

  // @ts-ignore
  return get(converter, style, Gradients.DefaultLight)
}

const Notifications = () => {
  const notifications = useSelector((state: AppState) => state.notifications)

  return (
    <>
      {notifications.map(item => {
        const {id, message, style, duration} = item
        const gradient = matchGradientToColor(style)

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
