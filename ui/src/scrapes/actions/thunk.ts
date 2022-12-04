import request from 'src/utils/request'
import {
  editScrape,
  removeScrape,
  addScrape,
  setScrapes,
} from 'src/scrapes/actions/creators'
import {defaultErrorNotification} from 'src/constants/notification'
import {notify} from 'src/shared/actions/notifications'
import {normalize} from 'normalizr'
import {ScrapeEntities} from 'src/types/schemas'
import {arrayOfScrapes, scrapeSchema} from 'src/schemas/scrapes'
import {Scrape} from 'src/types/scrape'
import {getByID} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {GetState} from 'src/types/stores'
import {getOrg} from 'src/organizations/selectors'
import {RemoteDataState} from '@influxdata/clockface'

export const getScrapes =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)

    try {
      const resp = await request(`/api/v1/scrapes?orgID=${org.id}`)
      if (resp.status !== 200) {
        throw new Error(resp.date.message)
      }

      const norm = normalize<Scrape, ScrapeEntities, string[]>(
        resp.data,
        arrayOfScrapes
      )

      dispatch(setScrapes(RemoteDataState.Done, norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Get scrapes failed, ${err}`,
        })
      )
    }
  }

export const deleteScrape =
  (id: string) =>
  async (dispatch): Promise<void> => {
    try {
      const resp = await request(`/api/v1/scrapes/${id}`, {method: 'DELETE'})
      if (resp.status !== 204) {
        throw new Error(resp.data.message)
      }

      dispatch(removeScrape(id))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Delete Scrape failed, ${err}`,
        })
      )
    }
  }

export interface ScrapeUpdate {
  name?: string
  desc?: string
}

export const updateScrape =
  (id: string, upd: ScrapeUpdate) =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const current = getByID<Scrape>(state, ResourceType.Scrapes, id)

    try {
      const resp = await request(`/api/v1/scrapes/${id}`, {
        method: 'PATCH',
        body: {
          name: upd.name || current.name,
          desc: upd.desc || current.desc,
        },
      })
      if (resp.status !== 200) {
        throw new Error(resp.data.message)
      }

      const normalized = normalize<Scrape, ScrapeEntities, string>(
        resp.data,
        scrapeSchema
      )

      dispatch(editScrape(normalized))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Update scrape failed, ${err}`,
        })
      )
    }
  }

export const createScrape =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)

    try {
      const resp = await request(`/api/v1/scrapes`, {
        method: 'POST',
        body: {
          name: 'selfstat',
          desc: 'Collect metrics of Manta',
          orgID: org.id,
          targets: ['localhost:8088'],
          labels: {
            foo: 'bar',
          },
        },
      })
      if (resp.status !== 201) {
        throw new Error(resp.data.message)
      }

      const norm = normalize<Scrape, ScrapeEntities, string>(
        resp.data,
        scrapeSchema
      )
      dispatch(addScrape(norm))
    } catch (err) {
      console.error(err)

      dispatch(
        notify({
          ...defaultErrorNotification,
          message: `Create scrape failed, ${err}`,
        })
      )
    }
  }
