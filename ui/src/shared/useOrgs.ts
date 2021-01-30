import constate from 'constate'
import {useFetch} from 'use-http'
import remoteDataState from '../utils/rds'

const [OrgsProvider, useOrgs] = constate(
  () => {
    const {data, error, loading} = useFetch(`/api/v1/orgs`, {}, [])

    return {
      error,
      orgs: data,
      loading: remoteDataState(data, error, loading),
    }
  },
  (values) => values
)

export {OrgsProvider, useOrgs}
