import constate from 'constate'
import {useCallback, useState} from 'react'

import {Notification} from 'src/types/Notification'

export * from './defaults'

const [NotificationProvider, useNotification] = constate(() => {
  const [notifications, setNotifications] = useState<Notification[]>([])

  const notify = useCallback((n: Notification) => {
    const id = `${Date.now() + performance.now()}`
    setNotifications(prev => [
      ...prev,
      {
        ...n,
        id,
      },
    ])
  }, [])

  const dismiss = useCallback((id?: string) => {
    setNotifications(prev => {
      return prev.filter(v => v.id !== id)
    })
  }, [])

  return {
    dismiss,
    notifications,
    notify,
  }
})

export {NotificationProvider, useNotification}
