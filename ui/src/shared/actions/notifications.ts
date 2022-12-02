import {Notification} from 'src/types/notification'

export const PUBLISH_NOTIFICATION = 'PUBLISH_NOTIFICATION'
export const DISMISS_NOTIFICATION = 'DISMISS_NOTIFICATION'

export type Action = ReturnType<typeof notify> | ReturnType<typeof dismiss>

export interface PublishNotificationAction {
  type: 'PUBLISH_NOTIFICATION'
  notification
}

export const notify = (notification: Notification) =>
  ({
    type: PUBLISH_NOTIFICATION,
    notification,
  } as const)

export const dismiss = (id: string) =>
  ({
    type: 'DISMISS_NOTIFICATION',
    id,
  } as const)
