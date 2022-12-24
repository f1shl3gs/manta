// Libraries
import {RemoteDataState} from '@influxdata/clockface'
import {produce} from 'immer'

// Types
import {ResourceState, ResourceType} from 'src/types/resources'
import {Check} from 'src/types/checks'

// Actions
import {
  Action,
  REMOVE_CHECK,
  SET_CHECK,
  SET_CHECKS,
  SET_CHECK_SEARCH_TERM,
} from 'src/checks/actions/creators'
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'

// Constants
import {DEFAULT_CHECK_SORT_OPTIONS} from 'src/shared/constants/checks'

type ChecksState = ResourceState['checks']

const initialState = (): ChecksState => ({
  byID: {},
  allIDs: [],
  searchTerm: '',
  status: RemoteDataState.NotStarted,
  sortOptions: DEFAULT_CHECK_SORT_OPTIONS,
})

export const checksReducer = (
  state: ChecksState = initialState(),
  action: Action
): ChecksState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_CHECKS:
        setResource<Check>(draftState, action, ResourceType.Checks)
        return

      case SET_CHECK:
        setResourceAtID<Check>(draftState, action, ResourceType.Checks)
        return

      case REMOVE_CHECK:
        removeResource<Check>(draftState, action)
        return

      case SET_CHECK_SEARCH_TERM:
        draftState.searchTerm = action.searchTerm
        return

      default:
        return
    }
  })
