import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'
import {RemoteDataState} from '@influxdata/clockface'
import {useOrgID} from '../../shared/useOrg'
import {Scraper} from '../../types/scrapers'
import {
  defaultDeletionNotification,
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from '../../shared/notification/useNotification'

const [ScrapersProvider, useScrapers] = constate(
  () => {
    const orgID = useOrgID()
    const {notify} = useNotification()
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const [list, setList] = useState<Scraper[]>([])
    const [reload, setReload] = useState(0)

    useEffect(() => {
      setLoading(RemoteDataState.Loading)
      fetch(`/api/v1/scrapes?orgID=${orgID}`)
        .then(resp => resp.json())
        .then(data => {
          setList(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Fetch scrapers failed, err: ${err.message}`,
          })
          setLoading(RemoteDataState.Error)
        })
    }, [notify, orgID, reload])

    const onRemove = useCallback(
      (id: string) => {
        fetch(`/api/v1/scrapes/${id}`, {method: 'DELETE'})
          .then(() => {
            notify({
              ...defaultDeletionNotification,
              message: `Delete scraper success`,
            })
            setReload(prevState => prevState + 1)
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Delete scraper failed, err: ${err.message}`,
            })
          })
      },
      [notify]
    )

    const onNameUpdate = useCallback(
      (id: string, name: string) => {
        const scraper = list.find(item => item.id === id)

        fetch(`/api/v1/scrapes/${id}`, {
          method: 'PATCH',
          body: JSON.stringify({
            name,
          }),
        })
          .then(() => {
            notify({
              ...defaultSuccessNotification,
              message: `Update scraper ${scraper?.name}'s name success`,
            })
            setReload(prev => prev + 1)
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Update scraper ${scraper?.name}'s name failed, err: ${err.message}`,
            })
          })
      },
      [list, notify]
    )

    const onDescUpdate = useCallback(
      (id: string, desc: string) => {
        const scraper = list.find(item => item.id === id)

        fetch(`/api/v1/scrapes/${id}`, {
          method: 'PATCH',
          body: JSON.stringify({desc}),
        })
          .then(() => {
            notify({
              ...defaultSuccessNotification,
              message: `Update scraper ${scraper?.name}'s desc success`,
            })
            setReload(prev => prev + 1)
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Update scraper ${scraper?.name}'s desc failed, err: ${err.message}`,
            })
          })
      },
      [list, notify]
    )

    return {
      loading,
      scrapers: list,
      onRemove,
      onNameUpdate,
      onDescUpdate,
    }
  },
  value => value
)

export {ScrapersProvider, useScrapers}
