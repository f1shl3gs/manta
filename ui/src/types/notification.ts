import { IconFont } from '@influxdata/clockface';

export enum NotificationStyle {
  Error = 'error',
  Success = 'success',
  Info = 'info',
  Primary = 'primary',
  Warning = 'warning',
}

export interface Notification {
  id?: string
  style: NotificationStyle
  icon: IconFont
  message: string
  duration?: number
  link?: string
  linkText?: string
}