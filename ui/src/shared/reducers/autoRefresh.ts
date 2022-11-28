// Libraries
import {produce} from 'immer'

// Actions
import {Action} from 'src/shared/actions/autoRefresh'

// Types
import {AutoRefresh, AutoRefreshStatus} from 'src/types/AutoRefresh'

export interface AutoRefreshState {
  autoRefresh: AutoRefresh
}

const initialState = (): AutoRefreshState => ({
  autoRefresh: {
    status: AutoRefreshStatus.Active,
    interval: 15,
  },
})

export const autoRefreshReducer = (state = initialState(), action: Action) =>
  produce(state, draftState => {
    switch (action.type) {
      case 'SetAutoRefreshInterval':
        draftState.autoRefresh.interval = action.payload.second
    }
  })
