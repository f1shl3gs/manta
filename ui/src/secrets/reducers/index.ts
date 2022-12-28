// Libraries
import {produce} from 'immer'

// Types
import {ResourceType, SecretsState} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {Secret} from 'src/types/secrets'

// Action
import {
  Action,
  REMOVE_SECRET,
  SET_SECRET,
  SET_SECRETS,
} from 'src/secrets/actions/creators'
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'

// Constants
import {DEFAULT_SECRET_SORT_OPTIONS} from 'src/secrets/constants'

const initialState = (): SecretsState => ({
  allIDs: [],
  byID: {},
  status: RemoteDataState.NotStarted,
  searchTerm: '',
  sortOptions: DEFAULT_SECRET_SORT_OPTIONS,
})

export const secretsReduer = (
  state: SecretsState = initialState(),
  action: Action
): SecretsState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_SECRETS:
        setResource<Secret>(draftState, action, ResourceType.Secrets)
        return

      case SET_SECRET:
        setResourceAtID<Secret>(draftState, action, ResourceType.Secrets)
        return

      case REMOVE_SECRET:
        removeResource(draftState, action)
        return

      default:
        return
    }
  })
