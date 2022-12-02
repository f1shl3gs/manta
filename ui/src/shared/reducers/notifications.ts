import {produce} from 'immer'
import {Action, PUBLISH_NOTIFICATION} from 'src/shared/actions/notifications'

import {Notification} from 'src/types/notification'

export const initialState: Notification[] = []

export const notificationsReducer = (
  state: Notification[] = initialState,
  action: Action
): Notification[] =>
  produce(state, draftState => {
    switch (action.type) {
      case PUBLISH_NOTIFICATION:
        const notification = {
          ...action.notification,
          id: `${Date.now() + performance.now()}`,
        }

        draftState.unshift(notification)
        return

      case 'DISMISS_NOTIFICATION':
        const {id} = action
        return draftState.filter(n => n.id !== id)

      default:
        return
    }
  })
