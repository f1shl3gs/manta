// Libraries
import {produce} from 'immer'

// Types
import {RemoteDataState, Sort} from '@influxdata/clockface'
import {NotificationEndpointState, ResourceType} from 'src/types/resources'
import {SortTypes} from 'src/types/sort'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'

// Actions
import {
  Action,
  REMOVE_NOTIFICATION_ENDPOINT,
  SET_NOTIFICATION_ENDPOINT,
  SET_NOTIFICATION_ENDPOINT_SEARCH_TERM,
  SET_NOTIFICATION_ENDPOINT_SORT_OPTIONS,
  SET_NOTIFICATION_ENDPOINTS,
} from 'src/notification_endpoints/actions/creators'

// Helper
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'

const initialState = (): NotificationEndpointState => ({
  allIDs: [],
  byID: {},
  status: RemoteDataState.NotStarted,
  searchTerm: '',
  sortOptions: {
    direction: Sort.Ascending,
    type: SortTypes.String,
    key: 'name',
  },
})

export const notificationEndpointsReducer = (
  state: NotificationEndpointState = initialState(),
  action: Action
): NotificationEndpointState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_NOTIFICATION_ENDPOINTS:
        setResource<NotificationEndpoint>(
          draftState,
          action,
          ResourceType.NotificationEndpoints
        )
        return

      case SET_NOTIFICATION_ENDPOINT:
        setResourceAtID<NotificationEndpoint>(
          draftState,
          action,
          ResourceType.NotificationEndpoints
        )
        return

      case REMOVE_NOTIFICATION_ENDPOINT:
        removeResource<NotificationEndpoint>(draftState, action)
        return

      case SET_NOTIFICATION_ENDPOINT_SEARCH_TERM:
        draftState.searchTerm = action.searchTerm
        return

      case SET_NOTIFICATION_ENDPOINT_SORT_OPTIONS:
        draftState.sortOptions = action.payload
        return

      default:
        return
    }
  })
