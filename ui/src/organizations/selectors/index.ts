// Libraries
import {get} from 'lodash'

// Types
import {AppState} from 'src/types/stores'
import {Organization} from 'src/types/Organization'
import {ResourceType} from 'src/types/resources'
import {getAll} from 'src/resources/selectors'
import {useSelector} from 'react-redux'

export const getOrg = (state: AppState): Organization => {
  return get(state, 'resources.orgs.org', null)
}

export const getOrgs = (state: AppState): Organization[] => {
  return getAll(state, ResourceType.Organizations)
}

export function useOrg() {
  return useSelector(getOrg)
}
