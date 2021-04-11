/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {Error} from '../models/Error'
import type {Otcl} from '../models/Otcl'
import type {OtclRequest} from '../models/OtclRequest'
import {request as __request} from '../core/request'

export class OtclsService {
  /**
   * List all OpenTelemetry Collector
   * @returns Otcl A list of Otcl
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async getOtcls({
    zapTraceSpan,
  }: {
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Otcl | Error> {
    const result = await __request({
      method: 'GET',
      path: `/otcls`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }

  /**
   * Create an Otcl
   * @returns Error Unexpected error
   * @returns Otcl Otcl created
   * @throws ApiError
   */
  public static async createOtcl({
    requestBody,
    zapTraceSpan,
  }: {
    /** Otcl to create **/
    requestBody: OtclRequest
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Error | Otcl> {
    const result = await __request({
      method: 'POST',
      path: `/otcls`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
      body: requestBody,
    })
    return result.body
  }

  /**
   * Get Otcl by ID
   * @returns Otcl Otcl details
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async getOtcl({
    id,
    zapTraceSpan,
  }: {
    /** The ID of the Otcl to get **/
    id: string
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Otcl | Error> {
    const result = await __request({
      method: 'GET',
      path: `/otcls/${id}`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }
}
