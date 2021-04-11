/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type OtclRequest = {
  name?: string
  desc?: string
  readonly orgID?: string
  type?: string
  /**
   * Config content for the OpenTelemetry Collector
   */
  content?: string
}
