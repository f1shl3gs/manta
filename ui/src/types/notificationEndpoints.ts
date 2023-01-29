import {RemoteDataState} from '@influxdata/clockface'

interface Base {
  readonly id?: string
  readonly created?: string
  readonly updated?: string

  name: string
  desc: string
  orgID: string

  type: string

  status: RemoteDataState
}

export type HTTPAuthMethod = 'none' | 'basic' | 'beaer'

export interface HTTP extends Base {
  method: string
  url: string
  headers: {
    [key: string]: string
  }

  authMethod: HTTPAuthMethod
  username?: string
  password?: string
  token?: string

  contentTemplate?: string
}

export type NotificationEndpoint = HTTP

export type NotificationEndpointType = NotificationEndpoint['type']
