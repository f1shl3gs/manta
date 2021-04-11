/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Otcl = {
  readonly id?: string
  readonly created?: string
  readonly modified?: string
  name?: string
  desc?: string
  readonly orgID?: string
  type?: string
  /**
   * Config content for the OpenTelemetry Collector
   */
  content?: string
}
