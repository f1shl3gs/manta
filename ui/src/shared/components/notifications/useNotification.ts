import constate from 'constate'
import {useCallback, useMemo, useState} from 'react'

import {Notification} from 'src/types/notification'

export * from 'src/shared/components/notifications/defaults'

const [NotificationProvider, useNotifications, useNotify] = constate(
  () => {
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
  },
  value =>
    useMemo(
      () => ({
        dismiss: value.dismiss,
        notifications: value.notifications,
      }),
      [value.dismiss, value.notifications]
    ),
  value => useMemo(() => value.notify, [value.notify])
)

export {NotificationProvider, useNotifications, useNotify}
