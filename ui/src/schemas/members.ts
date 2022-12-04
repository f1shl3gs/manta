import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

// Defines the schema for the 'members' resource
export const memberSchema = new schema.Entity(ResourceType.Members)
export const arrayOfMembers = [memberSchema]
