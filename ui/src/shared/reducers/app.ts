import {ActionTypes, Action} from 'src/shared/actions/app'

export interface AppState {
  inPresentationMode: boolean
}

const initialState: AppState = {
  inPresentationMode: false,
}

const appReducer = (state = initialState, action: Action): AppState => {
  switch (action.type) {
    case ActionTypes.DisablePresentationMode:
      return {
        ...state,
        inPresentationMode: false,
      }
    case ActionTypes.EnablePresentationMode:
      return {
        ...state,
        inPresentationMode: true,
      }

    default:
      return state
  }
}

export default appReducer
