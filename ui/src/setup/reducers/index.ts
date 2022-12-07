import {produce} from 'immer'

import {Action, SET_STEP, SET_SETUP_PARAMS} from 'src/setup/actions/creators'

export enum Step {
  Welcome,
  Admin,
}

interface SetupParams {
  username: string
  password: string
  organization: string
}

export interface SetupState extends SetupParams {
  step: Step
}

export const initialState = (): SetupState => ({
  step: Step.Welcome,
  username: '',
  password: '',
  organization: '',
})

export const setupReducer = (
  state: SetupState = initialState(),
  action: Action
) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_STEP:
        draftState.step = action.step
        return

      case SET_SETUP_PARAMS:
        const {username, password, organization} = action

        if (username) {
          draftState.username = username
        }

        if (password) {
          draftState.password = password
        }

        if (organization) {
          draftState.organization = organization
        }

        return

      default:
        return
    }
  })
