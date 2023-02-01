// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'

// Selectors
import {getResourcesStatus} from 'src/resources/selectors/GetResourcesStatus'
import {getResourceStatus} from 'src/resources/selectors/GetResourceStatus'

// Utils
import {get} from 'src/shared/utils/get'

export {getResourceStatus, getResourcesStatus}

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

export const getStatus = (
  {resources}: AppState,
  resource: ResourceType
): RemoteDataState => {
  return resources[resource].status
}
