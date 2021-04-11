/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type {Condition} from './Condition'

export type Check = {
  id?: string
  readonly created?: string
  readonly updated?: string
  name?: string
  desc?: string
  status?: Check.status
  orgID?: string
  expr?: string
  conditions?: Array<Condition>
}

export namespace Check {
  export enum status {
    ACTIVE = 'active',
    INACTIVE = 'inactive',
  }
}
