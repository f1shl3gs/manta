// Libraries
import constate from 'constate'
import {useCallback} from 'react'

// Hooks
import {useFetch} from 'shared/useFetch'

// Utils
import remoteDataState from 'utils/rds'

interface CheckUpdate {
  name?: string
  desc?: string
  status?: string
}

const [ChecksProvider, useChecks] = constate(
  () => {
    const {data, error, loading, get} = useFetch(`checks`, {}, [])
    const {patch} = useFetch(`checks`, {})
    const {del} = useFetch(`checks`, {})

    const patchCheck = useCallback(
      (id: string, udp: CheckUpdate) => {
        patch(id, udp)
          .then(() => {
            get()
          })
          .catch(err => {
            console.log('err', err)
          })
      },
      [get, patch]
    )

    return {
      patchCheck,
      remoteDataState: remoteDataState(data, error, loading),
      checks: data || [],
      reload: get,
      del: (id: string) => {
        return del(`/${id}`)
      },
    }
  },
  value => value
)

export {ChecksProvider, useChecks}
