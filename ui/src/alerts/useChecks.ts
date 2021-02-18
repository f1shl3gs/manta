// Libraries
import constate from 'constate'
import {useFetch} from 'use-http'

import {useOrgID} from '../shared/useOrg'
import remoteDataState from '../utils/rds'

const [ChecksProvider, useChecks] = constate(
  () => {
    const orgID = useOrgID()
    const {data, error, loading} = useFetch(
      `/api/v1/checks?orgID=${orgID}`,
      {},
      []
    )

    return {
      remoteDataState: remoteDataState(data, error, loading),
      checks: data || [],
    }
  },
  value => value
)

export {ChecksProvider, useChecks}
