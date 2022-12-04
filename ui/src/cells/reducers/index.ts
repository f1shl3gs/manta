import {produce} from 'immer'
import {get} from 'lodash'

import {ResourceState} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'

import {Action, REMOVE_CELL, SET_CELLS} from 'src/cells/actions/creators'

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

      case SET_CELLS:
        const {status, schema} = action
        draftState.status = status

        if (get(schema, ['entities', 'cells'])) {
          draftState.byID = schema.entities['cells']
        }
        return

      default:
        return
    }
  })