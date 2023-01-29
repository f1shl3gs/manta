// Types
import {NotificationEndpoint} from 'src/types/notificationEndpoints'
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'

export const getEndpoint = (state: AppState): NotificationEndpoint => {
  return state.resources[ResourceType.NotificationEndpoints].current
}
