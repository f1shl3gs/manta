// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'

import {getResourcesStatus} from './GetResourcesStatus'
import {getResourceStatus} from './GetResourceStatus'
import {RemoteDataState} from '@influxdata/clockface'

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
