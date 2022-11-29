// Libraries
import {produce} from 'immer'

// Types
import {OrgsState, ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {Organization} from 'src/types/Organization'

// Actions
import {Action, ADD_ORG, SET_ORG, SET_ORGS} from 'src/organizations/actions'

// Utils
import {addResource, setResource} from 'src/resources/reducers/helpers'

const initialState = (): OrgsState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
  org: null,
})

export const organizationsReducer = (
  state: OrgsState = initialState(),
  action: Action
) =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_ORGS:
        setResource<Organization>(
          draftState,
          action,
          ResourceType.Organizations
        )

        return
      case ADD_ORG:
        addResource<Organization>(
          draftState,
          action,
          ResourceType.Organizations
        )
        return
      case SET_ORG:
        draftState.org = action.org
        return
      default:
        return
    }
  })
