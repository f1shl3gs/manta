import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

export const scrapeSchema = new schema.Entity(ResourceType.Scrapes)

export const arrayOfScrapes = [scrapeSchema]
