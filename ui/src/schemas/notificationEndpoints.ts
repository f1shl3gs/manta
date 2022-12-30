import {schema} from 'normalizr'
import {ResourceType} from 'src/types/resources'

// Defines the schema for the 'NotificationEndpoint' resource
export const notificationEndpointSchema = new schema.Entity(
  ResourceType.NotificationEndpoints
)
export const arrayOfNotificationEndpoints = [notificationEndpointSchema]
