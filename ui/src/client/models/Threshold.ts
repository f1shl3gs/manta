/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Threshold = {
  type: Threshold.type
  value?: number
  min?: number
  max?: number
}

export namespace Threshold {
  export enum type {
    GT = 'gt',
    LT = 'lt',
    EQ = 'eq',
    INSIDE = 'inside',
    OUTSIDE = 'outside',
  }
}
