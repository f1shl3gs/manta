// Libraries
import constate from 'constate'
import {CachePolicies, useFetch} from 'shared/useFetch'

import {useOrgID} from '../../shared/useOrg'
import remoteDataState from '../../utils/rds'

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

    const {del} = useFetch(`/api/v1/checks`, {})

    return {
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
