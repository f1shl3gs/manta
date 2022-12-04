// Libraries
import {produce} from 'immer'

// Actions
import {
  Action,
  SET_AUTOREFRESH_INTERVAL,
  SET_RANGE,
} from 'src/shared/actions/autoRefresh'

// Types
import {AutoRefresh, AutoRefreshStatus} from 'src/types/autoRefresh'

export interface AutoRefreshState {
  autoRefresh: AutoRefresh

  start: number
  end: number
  step: number
}

const initialState = (): AutoRefreshState => ({
  autoRefresh: {
    status: AutoRefreshStatus.Active,
    interval: 15,
  },
  start: 0,
  end: 0,
  step: 0,
})

export const autoRefreshReducer = (state = initialState(), action: Action) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_AUTOREFRESH_INTERVAL:
        draftState.autoRefresh = action.payload
        return

      case SET_RANGE:
        const {start, end, step} = action.payload

        draftState.start = start
        draftState.end = end
        draftState.step = step

        return

      default:
        return
    }
  })
