/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {Check} from '../models/Check'
import type {Checks} from '../models/Checks'
import type {Error} from '../models/Error'
import {request as __request} from '../core/request'

export class ChecksService {
  /**
   * List al Checks
   * @returns Checks Check List
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async listChecks({
    zapTraceSpan,
    orgId,
  }: {
    /** OpenTracing span context **/
    zapTraceSpan?: string
    /** List checks with this orgID **/
    orgId?: string
  }): Promise<Checks | Error> {
    const result = await __request({
      method: 'GET',
      path: `/checks`,
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
   * Create a new Check
   * @returns Error Unexpected error
   * @returns any Create Check success
   * @throws ApiError
   */
  public static async createChecks({
    zapTraceSpan,
    requestBody,
  }: {
    /** OpenTracing span context **/
    zapTraceSpan?: string
    /** Check to create **/
    requestBody?: Check
  }): Promise<Error | any> {
    const result = await __request({
      method: 'PUT',
      path: `/checks`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
      body: requestBody,
    })
    return result.body
  }

  /**
   * Delete Check by id
   * @returns Error Unexpected error
   * @returns any Delete check by id success
   * @throws ApiError
   */
  public static async deleteCheck({
    id,
    zapTraceSpan,
  }: {
    id: string
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Error | any> {
    const result = await __request({
      method: 'DELETE',
      path: `/checks/${id}`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }
}
