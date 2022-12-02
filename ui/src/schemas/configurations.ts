import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

// Defines the schema for the 'Configurations' resource
export const configurationSchema = new schema.Entity(
  ResourceType.Configurations
)

export const arrayOfConfigurations = [configurationSchema]
