import {RemoteDataState} from '@influxdata/clockface'
import {NormalizedSchema} from 'normalizr'

import {DashboardEntities} from 'src/types/schemas'
import {DashboardSortParams} from 'src/types/dashboards'

// Action Types
export const EDIT_DASHBOARD = 'EDIT_DASHBOARD'
export const SET_DASHBOARDS = 'SET_DASHBOARDS'
export const REMOVE_DASHBOARD = 'REMOVE_DASHBOARD'
export const SET_DASHBOARD = 'SET_DASHBOARD'
export const SET_DASHBOARD_SORT = 'SET_DASHBOARD_SORT'
export const SET_DASHBOARD_SEARCH_TERM = 'SET_DASHBOARD_SEARCH_TERM'

export type Action =
  | ReturnType<typeof editDashboard>
  | ReturnType<typeof setDashboards>
  | ReturnType<typeof setDashboard>
  | ReturnType<typeof removeDashboard>
  | ReturnType<typeof setDashboardSort>
  | ReturnType<typeof setDashboardSearchTerm>

type DashboardSchema<R extends string | string[]> = NormalizedSchema<
  DashboardEntities,
  R
>

export const setDashboard = (
  id: string,
  status: RemoteDataState,
  schema?: DashboardSchema<string>
) =>
  ({
    type: SET_DASHBOARD,
    id,
    status,
    schema,
  } as const)

export const setDashboards = (
  status: RemoteDataState,
  schema?: DashboardSchema<string[]>
) =>
  ({
    type: SET_DASHBOARDS,
    status,
    schema,
  } as const)

export const removeDashboard = (id: string) =>
  ({
    type: REMOVE_DASHBOARD,
    id,
  } as const)

export const setDashboardSort = (sortOptions: DashboardSortParams) =>
  ({
    type: SET_DASHBOARD_SORT,
    sortOptions,
  } as const)

export const setDashboardSearchTerm = (searchTerm: string) =>
  ({
    type: SET_DASHBOARD_SEARCH_TERM,
    searchTerm,
  } as const)

export const editDashboard = (schema: DashboardSchema<string>) =>
  ({
    type: EDIT_DASHBOARD,
    schema,
  } as const)
