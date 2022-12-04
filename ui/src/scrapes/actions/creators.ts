import {RemoteDataState} from '@influxdata/clockface'
import {NormalizedSchema} from 'normalizr'
import {ScrapeEntities} from 'src/types/schemas'

export const SET_SCRAPES = 'SET_SCRAPES'
export const ADD_SCRAPE = 'SET_SCRAPE'
export const EDIT_SCRAPE = 'EDIT_SCRAPE'
export const REMOVE_SCRAPE = 'REMOVE_SCRAPE'

// R is the type of the value of the 'result' key in normalizr's normalization
type ScrapeSchema<R extends string | string[]> = NormalizedSchema<
  ScrapeEntities,
  R
>

export const setScrapes = (
  status: RemoteDataState,
  schema?: NormalizedSchema<ScrapeEntities, string[]>
) =>
  ({
    type: SET_SCRAPES,
    status,
    schema,
  } as const)

export const addScrape = (schema?: ScrapeSchema<string>) =>
  ({
    type: ADD_SCRAPE,
    schema,
  } as const)

export const editScrape = (schema: NormalizedSchema<ScrapeEntities, string>) =>
  ({
    type: EDIT_SCRAPE,
    schema,
  } as const)

export const removeScrape = (id: string) =>
  ({
    type: REMOVE_SCRAPE,
    id,
  } as const)

export type Action =
  | ReturnType<typeof setScrapes>
  | ReturnType<typeof addScrape>
  | ReturnType<typeof editScrape>
  | ReturnType<typeof removeScrape>
