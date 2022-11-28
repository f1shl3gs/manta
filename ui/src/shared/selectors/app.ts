import {AppState} from 'src/types/stores'

export const getPresentationMode = (state: AppState): boolean =>
  state.app.inPresentationMode
