import {RemoteDataState} from '@influxdata/clockface'
import produce from 'immer'
import {MembersState, ResourceType} from 'src/types/resources'
import {Action, SET_MEMBERS} from 'src/members/actions'
import {setResource} from 'src/resources/reducers/helpers'

const initialState = (): MembersState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  me: null,
})

export const usersReducer = (
  state: MembersState = initialState(),
  action: Action
) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_MEMBERS:
        setResource(draftState, action, ResourceType.Users)
        return
      default:
        return
    }
  })
