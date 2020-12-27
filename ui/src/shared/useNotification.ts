import { useCallback, useState } from 'react';
import constate from 'constate';

import { Notification } from 'types/notification';

const [NotificationProvider, useNotification] = constate(
  () => {
    const [notifications, setNotifications] = useState<Notification[]>([]);

    const info = useCallback((msg: string) => {

    }, [notifications]);

    return {
      notifications,
      notification: {
        info
      }
    };
  }
);

export {
  NotificationProvider,
  useNotification
};