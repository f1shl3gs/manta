import {AppState} from 'src/types/stores'

export const getPresentationMode = (state: AppState): boolean =>
  state.app.ephemeral.inPresentationMode

export const getNavBarState = (state: AppState) =>
  state.app.persisted.navbarState
