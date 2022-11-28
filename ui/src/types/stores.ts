import {AppState as State} from 'src/shared/reducers/app'
// import {Notification} from 'src/types/Notification'

export interface AppState {
  app: State

  // notifications: Notification[]
}

export type GetState = () => AppState
