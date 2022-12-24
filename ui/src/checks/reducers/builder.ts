// Libraries
import {produce} from 'immer'
import {RemoteDataState} from '@influxdata/clockface'

// Types
import {CheckStatus, Conditions} from 'src/types/checks'

// Constants
import {
  DEFAULT_CHECK_DESC,
  DEFAULT_CHECK_CRON,
  DEFAULT_CHECK_NAME,
} from 'src/checks/constants'

// Actions
import {
  Action,
  REMOVE_CHECK_CONDITION,
  RESET_CHECK_BUILDER,
  SET_CHECK_BUILDER,
  SET_CHECK_CONDITIONS,
  SET_CHECK_CRON,
  SET_CHECK_NAME,
  SET_CHECK_QUERY,
  SET_CHECK_BUILDER_TAB,
  SET_CHECK_CONDITION,
} from 'src/checks/actions/builder'

export interface CheckBuilderState {
  readonly id: string

  name: string
  desc: string
  activeStatus: CheckStatus
  conditions: Conditions
  cron: string
  query: string
  tab: 'query' | 'meta'

  status: RemoteDataState
}

const initialState = (): CheckBuilderState => ({
  id: null,
  name: DEFAULT_CHECK_NAME,
  desc: DEFAULT_CHECK_DESC,
  activeStatus: 'active',
  conditions: {},
  cron: DEFAULT_CHECK_CRON,
  query: '',
  tab: 'query',
  status: RemoteDataState.NotStarted,
})

export const checkBuilderReducer = (
  state: CheckBuilderState = initialState(),
  action: Action
): CheckBuilderState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_CHECK_NAME:
        draftState.name = action.payload.name
        return

      case SET_CHECK_CRON:
        draftState.cron = action.payload.cron
        return

      case SET_CHECK_QUERY:
        draftState.query = action.payload.query
        return

      case SET_CHECK_CONDITIONS:
        draftState.conditions = action.payload.conditions
        return

      case SET_CHECK_BUILDER_TAB:
        draftState.tab = action.payload.tab
        return

      case RESET_CHECK_BUILDER:
        return initialState()

      case SET_CHECK_BUILDER:
        const {status, id, name, query, desc, conditions, cron} = action.payload

        return {
          ...draftState,
          id,
          name,
          desc,
          status,
          conditions,
          cron,
          query,
        }

      case SET_CHECK_CONDITION:
        const condition = action.payload
        draftState.conditions[condition.status] = condition
        return

      case REMOVE_CHECK_CONDITION:
        draftState.conditions[action.payload.status] = undefined
        return

      default:
        return
    }
  })
