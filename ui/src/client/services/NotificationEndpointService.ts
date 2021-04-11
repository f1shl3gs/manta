/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {Error} from '../models/Error'
import type {NotificationEndpoint} from '../models/NotificationEndpoint'
import type {NotificationEndpoints} from '../models/NotificationEndpoints'
import {request as __request} from '../core/request'

export class NotificationEndpointService {
  /**
   * List all notification
   * @returns NotificationEndpoints Notification Endpoint List
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async listNotificationEndpoint({
    orgId,
    zapTraceSpan,
  }: {
    orgId: any
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<NotificationEndpoints | Error> {
    const result = await __request({
      method: 'GET',
      path: `/notification_endpoint`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
      query: {
        orgID: orgId,
      },
    })
    return result.body
  }

  /**
   * Get NotificationEndpoint
   * @returns NotificationEndpoint NotificationEndpoint Detail
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async getNotificationEndpoint({
    id,
    zapTraceSpan,
  }: {
    id: string
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<NotificationEndpoint | Error> {
    const result = await __request({
      method: 'GET',
      path: `/notification_endpoint/${id}`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }

  /**
   * Delete notification endpoint by id
   * @returns Error Unexpected error
   * @returns any Delete success
   * @throws ApiError
   */
  public static async deleteNotificationEndpoint({
    id,
    zapTraceSpan,
  }: {
    id: string
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Error | any> {
    const result = await __request({
      method: 'DELETE',
      path: `/notification_endpoint/${id}`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }
}
