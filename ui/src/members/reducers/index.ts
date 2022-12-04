import {RemoteDataState} from '@influxdata/clockface'
import produce from 'immer'
import {UsersState, ResourceType} from 'src/types/resources'
import {Action, SET_MEMBERS} from 'src/members/actions/creators'
import {setResource} from 'src/resources/reducers/helpers'

const initialState = (): UsersState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  me: null,
})

export const usersReducer = (
  state: UsersState = initialState(),
  action: Action
) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_MEMBERS:
        setResource(draftState, action, ResourceType.Members)
        return

      default:
        return
    }
  })
