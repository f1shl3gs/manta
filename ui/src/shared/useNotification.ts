import constate from 'constate'
import {useCallback, useState} from 'react'

import {Notification} from 'types/Notification'

const [NotificationProvider, useNotification] = constate(() => {
  const [notifications, setNotifications] = useState<Notification[]>([])

  const info = useCallback((msg: string) => {}, [notifications])

  return {
    notifications,
    notification: {
      info,
    },
  }
})

export {NotificationProvider, useNotification}
