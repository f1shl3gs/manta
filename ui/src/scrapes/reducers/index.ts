import {produce} from 'immer'

import {ResourceType, ScrapesState} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {Action, REMOVE_SCRAPE, SET_SCRAPES} from 'src/scrapes/actions/creators'
import {removeResource, setResource} from 'src/resources/reducers/helpers'

const initialState = (): ScrapesState => ({
  byID: {},
  allIDs: [],
  status: RemoteDataState.NotStarted,
})

export const scrapesReducers = (
  state: ScrapesState = initialState(),
  action: Action
): ScrapesState =>
  produce(state, draftState => {
    switch (action.type) {
      case SET_SCRAPES:
        setResource(draftState, action, ResourceType.Scrapes)
        return
      case REMOVE_SCRAPE:
        removeResource(draftState, action)
        return
      default:
        return
    }
  })
