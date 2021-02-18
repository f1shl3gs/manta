import constate from 'constate'

const [NotificationEndpointsProvider, useNotificationEndpoints] = constate(
  () => {},
  value => value
)

export {NotificationEndpointsProvider, useNotificationEndpoints}
