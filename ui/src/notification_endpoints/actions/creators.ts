// Libraries
import {NormalizedSchema} from 'normalizr'

// Types
import {RemoteDataState, Sort} from '@influxdata/clockface'
import {NotificationEndpointEntities} from 'src/types/schemas'
import {SortTypes} from 'src/types/sort'

// Action types
export const SET_NOTIFICATION_ENDPOINTS = 'SET_NOTIFICATION_ENDPOINTS'
export const SET_NOTIFICATION_ENDPOINT = 'SET_NOTIFICATION_ENDPOINT'
export const REMOVE_NOTIFICATION_ENDPOINT = 'REMOVE_NOTIFICATION_ENDPOINT'
export const SET_NOTIFICATION_ENDPOINT_SEARCH_TERM =
  'SET_NOTIFICATION_ENDPOINT_SEARCH_TERM'
export const SET_NOTIFICATION_ENDPOINT_SORT_OPTIONS =
  'SET_NOTIFICATION_ENDPOINT_SORT_OPTIONS'

type NotificationEndpointSchema<R extends string | string[]> = NormalizedSchema<
  NotificationEndpointEntities,
  R
>

export const setNotificationEndpoints = (
  status: RemoteDataState,
  schema?: NotificationEndpointSchema<string[]>
) =>
  ({
    type: SET_NOTIFICATION_ENDPOINTS,
    status,
    schema,
  } as const)

export const setNotificationEndpoint = (
  id: string,
  status: RemoteDataState,
  schema?: NotificationEndpointSchema<string>
) =>
  ({
    type: SET_NOTIFICATION_ENDPOINT,
    id,
    status,
    schema,
  } as const)

export const removeNotificationEndpoint = (id: string) =>
  ({
    type: REMOVE_NOTIFICATION_ENDPOINT,
    id,
  } as const)

export const setNotificationEndpointSearchTerm = (searchTerm: string) =>
  ({
    type: SET_NOTIFICATION_ENDPOINT_SEARCH_TERM,
    searchTerm,
  } as const)

export const setNotificationEndpointSortOptions = (
  key: string,
  direction: Sort,
  type: SortTypes
) =>
  ({
    type: SET_NOTIFICATION_ENDPOINT_SORT_OPTIONS,
    payload: {
      key,
      direction,
      type,
    },
  } as const)

export type Action =
  | ReturnType<typeof setNotificationEndpoints>
  | ReturnType<typeof setNotificationEndpoint>
  | ReturnType<typeof removeNotificationEndpoint>
  | ReturnType<typeof setNotificationEndpointSearchTerm>
  | ReturnType<typeof setNotificationEndpointSortOptions>
