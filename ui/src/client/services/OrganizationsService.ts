/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {Error} from '../models/Error'
import type {Organization} from '../models/Organization'
import type {Organizations} from '../models/Organizations'
import {request as __request} from '../core/request'

export class OrganizationsService {
  /**
   * List all organizations
   * @returns Organizations A list of organizations
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async getOrgs({
    zapTraceSpan,
  }: {
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Organizations | Error> {
    const result = await __request({
      method: 'GET',
      path: `/orgs`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }

  /**
   * Create an Organization
   * @returns Organization Unexpected error
   * @throws ApiError
   */
  public static async createOrg({
    requestBody,
    zapTraceSpan,
  }: {
    /** Organization to create **/
    requestBody: Organization
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Organization> {
    const result = await __request({
      method: 'POST',
      path: `/orgs`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
      body: requestBody,
    })
    return result.body
  }

  /**
   * Retrieve an organization
   * @returns Organization Organization details
   * @returns Error Unexpected error
   * @throws ApiError
   */
  public static async getOrgById({
    id,
    zapTraceSpan,
  }: {
    /** The ID of the organization to get **/
    id: string
    /** OpenTracing span context **/
    zapTraceSpan?: string
  }): Promise<Organization | Error> {
    const result = await __request({
      method: 'GET',
      path: `/orgs/${id}`,
      headers: {
        'Zap-Trace-Span': zapTraceSpan,
      },
    })
    return result.body
  }
}
