import {ViewProperties, ViewType} from 'src/types/cells'
import {RemoteDataState} from '@influxdata/clockface'
import {FromFluxResult} from '@influxdata/giraffe'

export const SET_ACTIVE_QUERY = 'SET_ACTIVE_QUERY'
export const SET_ACTIVE_QUERY_TEXT = 'SET_ACTIVE_QUERY_TEXT'
export const ADD_QUERY = 'ADD_QUERY'
export const REMOVE_QUERY = 'REMOVE_QUERY'
export const SET_VIEWING_VIS_OPTIONS = 'SET_VIEWING_VIS_OPTIONS'
export const SET_VIEWTYPE = 'SET_VIEWTYPE'
export const SET_VIEW_PROPERTIES = 'SET_VIEW_PROPERTIES'
export const SET_VIEW_NAME = 'SET_VIEW_NAME'
export const RESET_TIMEMACHINE = 'RESET_TIMEMACHINE'
export const SET_QUERY_RESULTS = 'SET_QUERY_RESULTS'

export const setActiveQueryIndex = (activeQueryIndex: number) =>
  ({
    type: SET_ACTIVE_QUERY,
    activeQueryIndex,
  } as const)

export const setActiveQueryText = (text: string) =>
  ({
    type: SET_ACTIVE_QUERY_TEXT,
    text,
  } as const)

export const addQuery = () =>
  ({
    type: ADD_QUERY,
    query: {
      name: '',
      text: '',
      hidden: false,
    },
  } as const)

export const removeQuery = (index: number) =>
  ({
    type: REMOVE_QUERY,
    index,
  } as const)

export const setViewingVisOptions = (viewingVisOptions: boolean) =>
  ({
    type: SET_VIEWING_VIS_OPTIONS,
    viewingVisOptions,
  } as const)

export const setViewType = (viewType: ViewType) =>
  ({
    type: SET_VIEWTYPE,
    viewType,
  } as const)

export const setViewProperties = (viewProperties: ViewProperties) =>
  ({
    type: SET_VIEW_PROPERTIES,
    viewProperties,
  } as const)

export const setViewName = (name: string) =>
  ({
    type: SET_VIEW_NAME,
    name,
  } as const)

export const resetTimeMachine = () =>
  ({
    type: RESET_TIMEMACHINE,
  } as const)

export const setQueryResult = (
  state: RemoteDataState,
  result?: FromFluxResult
) =>
  ({
    type: SET_QUERY_RESULTS,
    payload: {
      state,
      result,
    },
  } as const)

export type Action =
  | ReturnType<typeof setActiveQueryIndex>
  | ReturnType<typeof setActiveQueryText>
  | ReturnType<typeof addQuery>
  | ReturnType<typeof removeQuery>
  | ReturnType<typeof setViewingVisOptions>
  | ReturnType<typeof setViewType>
  | ReturnType<typeof setViewProperties>
  | ReturnType<typeof setViewName>
  | ReturnType<typeof resetTimeMachine>
  | ReturnType<typeof setQueryResult>
