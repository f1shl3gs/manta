import {Notification} from 'src/types/notification'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
} from 'src/shared/constants/notification'

export const PUBLISH_NOTIFICATION = 'PUBLISH_NOTIFICATION'
export const DISMISS_NOTIFICATION = 'DISMISS_NOTIFICATION'

export type Action =
  | ReturnType<typeof notify>
  | ReturnType<typeof dismissNotification>

export interface PublishNotificationAction {
  type: 'PUBLISH_NOTIFICATION'
  notification
}

export const notify = (notification: Notification) =>
  ({
    type: PUBLISH_NOTIFICATION,
    notification,
  } as const)

export const dismissNotification = (id: string) =>
  ({
    type: DISMISS_NOTIFICATION,
    id,
  } as const)

export const error = (message: string) => ({
  type: PUBLISH_NOTIFICATION,
  notification: {
    ...defaultErrorNotification,
    message,
  },
})

export const info = (message: string) => ({
  type: PUBLISH_NOTIFICATION,
  notification: {
    ...defaultSuccessNotification,
    message,
  },
})
