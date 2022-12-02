import {AppState} from 'src/shared/reducers/app'
import {ResourceState} from 'src/types/resources'

export interface LocalStorage {
  app: AppState
  resources: {
    orgs: ResourceState['organizations']
  }
}
