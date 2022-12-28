// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {Secret} from 'src/types/secrets'

// Helper
import {getAll} from 'src/resources/selectors'

export const getAllSecrets = (state: AppState) => {
  return getAll<Secret>(state, ResourceType.Secrets)
}
