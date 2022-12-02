import {produce} from 'immer'

import {ResourceState} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'

import {Action, REMOVE_CELL} from 'src/cells/actions/creators'

type CellsState = ResourceState['cells']
const initialState = (): CellsState => ({
  byID: {},
  status: RemoteDataState.NotStarted,
})

export const cellsReducer = (
  state: CellsState = initialState(),
  action: Action
) =>
  produce(state, draftState => {
    switch (action.type) {
      case REMOVE_CELL:
        delete draftState.byID[action.id]
        return
      default:
        return
    }
  })
