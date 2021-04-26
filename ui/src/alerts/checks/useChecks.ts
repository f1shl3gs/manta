// Libraries
import constate from 'constate'
import {CachePolicies, useFetch} from 'shared/useFetch'

import {useOrgID} from 'shared/useOrg'
import remoteDataState from 'utils/rds'
import {useCallback} from 'react'

interface CheckUpdate {
  name?: string
  desc?: string
  status?: string
}

const [ChecksProvider, useChecks] = constate(
  () => {
    const orgID = useOrgID()
    const {data, error, loading, get} = useFetch(
      `/api/v1/checks?orgID=${orgID}`,
      {
        cachePolicy: CachePolicies.NO_CACHE,
      },
      []
    )

    const {patch} = useFetch(`/api/v1/checks/`, {})

    const {del} = useFetch(`/api/v1/checks`, {})

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
        return del(`/${id}?orgID=${orgID}`)
      },
    }
  },
  value => value
)

export {ChecksProvider, useChecks}
