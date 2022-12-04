import request from 'src/utils/request'
import {editScrape, removeScrape} from 'src/scrapes/actions/creators'
import {defaultErrorNotification} from 'src/constants/notification'
import {notify} from 'src/shared/actions/notifications'
import {normalize} from 'normalizr'
import {ScrapeEntities} from 'src/types/schemas'
import {scrapeSchema} from 'src/schemas/scrapes'
import {Scrape} from 'src/types/scrape'
import {getByID} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {GetState} from 'src/types/stores'

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
