import {ActionTypes, Action} from 'src/shared/actions/app'
import {NavBarState} from 'src/types/app'
import {combineReducers} from 'redux'

export interface AppState {
  ephemeral: {
    inPresentationMode: boolean
  }

  persisted: {
    navbarState: NavBarState
  }
}

const initialState: AppState = {
  ephemeral: {
    inPresentationMode: false
  },

  persisted: {
    navbarState: 'collapsed'
  }
}

const ephemeralReducer = (
  state = initialState.ephemeral,
    action: Action
): AppState['ephemeral'] => {
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

const persistedReducer = (
  state = initialState.persisted,
    action: Action
): AppState['persisted'] => {
  switch (action.type) {
    case ActionTypes.ToggleNavBarState:
      const navbarState = state.navbarState === 'expanded' ? 'collapsed' : 'expanded'

      return {
        ...state,
        navbarState
      }

  default:
    return state
  }
}

const appReducer = combineReducers<AppState>({
  ephemeral: ephemeralReducer,
  persisted: persistedReducer
})

export default appReducer
