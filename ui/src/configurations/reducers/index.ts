import produce from 'immer'
import {ConfigurationsState, ResourceType} from 'src/types/resources'
import {
  Action,
  REMOVE_CONFIG,
  SET_CONFIG,
  SET_CONFIGS,
} from 'src/configurations/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'
import {Configuration} from 'src/types/configuration'
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'

const initialState = (): ConfigurationsState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  config: {
    status: RemoteDataState.NotStarted,
    content: '',
  },
})

export const configurationsReducer = (
  state: ConfigurationsState = initialState(),
  action: Action
): ConfigurationsState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_CONFIGS:
        setResource<Configuration>(
          draftState,
          action,
          ResourceType.Configurations
        )
        return

      case SET_CONFIG:
        setResourceAtID<Configuration>(
          draftState,
          action,
          ResourceType.Configurations
        )
        return

      case REMOVE_CONFIG:
        removeResource<Configuration>(draftState, action)
        return

      default:
        return
    }
  })
