import React from 'react'
import {get} from 'lodash'

import {
  ComponentSize,
  Gradients,
  IconFont,
  Notification,
} from '@influxdata/clockface'

import {useNotifications} from 'src/shared/components/notifications/useNotification'
import {NotificationStyle} from 'src/types/Notification'

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
  const {notifications, dismiss} = useNotifications()

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
