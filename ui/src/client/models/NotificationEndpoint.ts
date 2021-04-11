/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type NotificationEndpoint = {
  readonly id: string
  readonly created: string
  readonly updated: string
  name: string
  desc: string
  status?: NotificationEndpoint.status
  url?: string
  method?: NotificationEndpoint.method
  headers?: any
  content: string
}

export namespace NotificationEndpoint {
  export enum status {
    ACTIVE = 'active',
    INACTIVE = 'inactive',
  }

  export enum method {
    POST = 'POST',
    GET = 'GET',
    PUT = 'PUT',
  }
}
