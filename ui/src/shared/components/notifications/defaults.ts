import {NotificationStyle} from 'types/Notification'
import {IconFont} from '@influxdata/clockface'
import {Notification} from 'types/Notification'

export const FIVE_SECONDS = 5000
export const TEN_SECONDS = 10000
export const FIFTEEN_SECONDS = 15000

type NotificationExcludingMessage = Pick<
  Notification,
  Exclude<keyof Notification, 'message'>
>

export const defaultErrorNotification: NotificationExcludingMessage = {
  style: NotificationStyle.Error,
  icon: IconFont.AlertTriangle,
  duration: TEN_SECONDS,
}

export const defaultSuccessNotification: NotificationExcludingMessage = {
  style: NotificationStyle.Success,
  icon: IconFont.CheckMark_New,
  duration: FIVE_SECONDS,
}

export const defaultDeletionNotification: NotificationExcludingMessage = {
  style: NotificationStyle.Primary,
  icon: IconFont.Trash_New,
  duration: FIVE_SECONDS,
}
