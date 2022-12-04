import {RemoteDataState} from '@influxdata/clockface'
import {NormalizedSchema} from 'normalizr'
import {ScrapeEntities} from 'src/types/schemas'

export const SET_SCRAPES = 'SET_SCRAPES'
export const EDIT_SCRAPE = 'EDIT_SCRAPE'
export const REMOVE_SCRAPE = 'REMOVE_SCRAPE'

export const setScrapes = (
  status: RemoteDataState,
  schema?: NormalizedSchema<ScrapeEntities, string[]>
) =>
  ({
    type: SET_SCRAPES,
    status,
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
  | ReturnType<typeof editScrape>
  | ReturnType<typeof removeScrape>
