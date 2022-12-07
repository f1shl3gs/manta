import {produce} from 'immer'

import {Action, SET_ME} from 'src/me/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'

export interface MeState {
  id: string
  name: string
  state: RemoteDataState
}

const initialState = (): MeState => ({
  id: '',
  name: '',
  state: RemoteDataState.NotStarted,
})

export const meReducer = (state: MeState = initialState(), action: Action) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_ME:
        const {id, name, state} = action

        draftState.id = id
        draftState.name = name
        draftState.state = state

        return

      default:
        return
    }
  })
