import produce from 'immer'
import {ConfigsState, ResourceType} from 'src/types/resources'
import {
  Action,
  REMOVE_CONFIG,
  SET_CONFIG,
  SET_CONFIGS,
} from 'src/configs/actions/creators'
import {RemoteDataState} from '@influxdata/clockface'
import {Config} from 'src/types/config'
import {
  removeResource,
  setResource,
  setResourceAtID,
} from 'src/resources/reducers/helpers'

const initialState = (): ConfigsState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  config: {
    status: RemoteDataState.NotStarted,
    content: '',
  },
})

export const configsReducer = (
  state: ConfigsState = initialState(),
  action: Action
): ConfigsState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_CONFIGS:
        setResource<Config>(draftState, action, ResourceType.Configs)
        return

      case SET_CONFIG:
        setResourceAtID<Config>(draftState, action, ResourceType.Configs)
        return

      case REMOVE_CONFIG:
        removeResource<Config>(draftState, action)
        return

      default:
        return
    }
  })
