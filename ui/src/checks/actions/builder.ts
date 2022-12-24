// Types
import {Check, Condition, Conditions, ConditionStatus} from 'src/types/checks'

// Action Types
export const SET_CHECK_NAME = 'SET_CHECK_NAME'
export const SET_CHECK_CRON = 'SET_CHECK_CRON'
export const SET_CHECK_QUERY = 'SET_CHECK_QUERY'
export const SET_CHECK_CONDITIONS = 'SET_CHECK_CONDITIONS'
export const SET_CHECK_BUILDER_TAB = 'SET_CHECK_BUILDER_TAB'
export const RESET_CHECK_BUILDER = 'RESET_CHECK_BUILDER'
export const SET_CHECK_BUILDER = 'SET_CHECK_BUILDER'
export const SET_CHECK_CONDITION = 'SET_CHECK_CONDITION'
export const REMOVE_CHECK_CONDITION = 'REMOVE_CHECK_CONDITION'

export const setName = (name: string) =>
  ({
    type: SET_CHECK_NAME,
    payload: {name},
  } as const)

export const setCron = (cron: string) =>
  ({
    type: SET_CHECK_CRON,
    payload: {
      cron,
    },
  } as const)

export const setQuery = (query: string) =>
  ({
    type: SET_CHECK_QUERY,
    payload: {
      query,
    },
  } as const)

export const setConditions = (conditions: Conditions) =>
  ({
    type: SET_CHECK_CONDITIONS,
    payload: {
      conditions,
    },
  } as const)

export const resetCheckBuilder = () =>
  ({
    type: RESET_CHECK_BUILDER,
  } as const)

export const setCheckBuilder = (check: Check) =>
  ({
    type: SET_CHECK_BUILDER,
    payload: check,
  } as const)

export const setTab = (tab: 'query' | 'meta') =>
  ({
    type: SET_CHECK_BUILDER_TAB,
    payload: {tab},
  } as const)

export const setCondition = (condition: Condition) =>
  ({
    type: SET_CHECK_CONDITION,
    payload: condition,
  } as const)

export const removeCondition = (status: ConditionStatus) =>
  ({
    type: REMOVE_CHECK_CONDITION,
    payload: {
      status,
    },
  } as const)

export type Action =
  | ReturnType<typeof setName>
  | ReturnType<typeof setCron>
  | ReturnType<typeof setQuery>
  | ReturnType<typeof setConditions>
  | ReturnType<typeof resetCheckBuilder>
  | ReturnType<typeof setCheckBuilder>
  | ReturnType<typeof setTab>
  | ReturnType<typeof setCondition>
  | ReturnType<typeof removeCondition>
