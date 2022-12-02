// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'

export const getByID = <R>(
  {resources}: AppState,
  type: ResourceType,
  id: string
): R => {
  const byID = get(resources, `${type}.byID`)

  if (!byID) {
    throw new Error(`"${type}" resource has yet not been set`)
  }

  return get(byID, `${id}`, null)
}

export const getAll = <R>(
  {resources}: AppState,
  resource: ResourceType
): R[] => {
  const allIDs: string[] = resources[resource].allIDs
  const byID: {[uuid: string]: R} = resources[resource].byID

  return (allIDs ?? []).map(id => byID[id])
}

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

export const getStatus = (
  {resources}: AppState,
  resource: ResourceType
): RemoteDataState => {
  return resources[resource].status
}
