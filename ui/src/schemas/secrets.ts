import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

export const secretSchema = new schema.Entity(
  ResourceType.Secrets,
  {},
  {idAttribute: 'key'}
)

export const arrayOfSecrets = [secretSchema]
