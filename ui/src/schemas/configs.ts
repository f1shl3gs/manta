import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

// Defines the schema for the 'Configs' resource
export const configSchema = new schema.Entity(ResourceType.Configs)

export const arrayOfConfigs = [configSchema]
