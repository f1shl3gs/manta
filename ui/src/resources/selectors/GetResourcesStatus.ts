import {RemoteDataState} from '@influxdata/clockface'
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

export const getResourcesStatus = (
  state: AppState,
  resources: Array<ResourceType>
): RemoteDataState => {
  const statuses = resources.map(resource => {
    if (!state.resources || !state.resources[resource].status) {
      throw new Error(
        `RemoteDataState status for resource "${resource}" is undefined in getResourceStatus`
      )
    }

    return state.resources[resource].status
  })

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

export default getResourcesStatus
