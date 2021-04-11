/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {Error} from '../models/Error'
import type {OnBoardingRequest} from '../models/OnBoardingRequest'
import {request as __request} from '../core/request'

export class SetupService {
  /**
   * Setup the service
   * @returns Error Unexpected error
   * @returns any Initial success
   * @throws ApiError
   */
  public static async setup({
    requestBody,
    zapTraceSpan,
  }: {
    /** User and Organization to create **/
    requestBody: OnBoardingRequest
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Error | any> {
    const result = await __request({
      method: 'POST',
      path: `/setup`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
      body: requestBody,
    })
    return result.body
  }
}
