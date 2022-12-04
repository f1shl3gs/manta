import {get} from 'lodash'

import {RemoteDataState} from '@influxdata/clockface'
import {Resource} from 'src/types/resources'
import {AppState} from 'src/types/stores'

const getStatus = ({resources}: AppState, {type, id}: Resource) => {
  return get(resources, [type, 'byID', id, 'status'], RemoteDataState.Loading)
}

export const getResourceStatus = (
  state: AppState,
  resources: Resource[]
): RemoteDataState => {
  const statuses = resources.map(resource => getStatus(state, resource))

  let status = RemoteDataState.NotStarted
  if (statuses.every(s => s === RemoteDataState.Done)) {
    status = RemoteDataState.Done
  } else if (statuses.includes(RemoteDataState.Error)) {
    status = RemoteDataState.Error
  } else if (statuses.includes(RemoteDataState.Loading)) {
    status = RemoteDataState.Loading
  }

  return status
}
