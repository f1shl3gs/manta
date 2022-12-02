import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

// Defines the schema for the 'organizations' resource
export const orgSchema = new schema.Entity(ResourceType.Organizations)
export const arrayOfOrgs = [orgSchema]
