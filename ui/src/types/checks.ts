// Types
import {RemoteDataState, Sort} from '@influxdata/clockface'
import {SortTypes} from 'src/types/sort'

export interface GreatThanThreshold {
  type: 'gt'
  value: number
}

export interface GreatEqualThreshold {
  type: 'ge'
  value: number
}

export interface EqualThreshold {
  type: 'eq'
  value: number
}

export interface NotEqualThreshold {
  type: 'ne'
  value: number
}

export interface LessThanThreshold {
  type: 'lt'
  value: number
}

export interface LessEqualThreshold {
  type: 'le'
  value: number
}

interface Range {
  min: number
  max: number
}

export interface InsideThreshold extends Range {
  type: 'inside'
}

export interface OutsideThreshold extends Range {
  type: 'outside'
}

export type Threshold =
  | InsideThreshold
  | OutsideThreshold
  | GreatThanThreshold
  | GreatEqualThreshold
  | EqualThreshold
  | NotEqualThreshold
  | LessEqualThreshold
  | LessThanThreshold

export type ThresholdType = Threshold['type']

export type ConditionStatus = 'crit' | 'warn' | 'info' | 'ok'

export interface Condition {
  status: ConditionStatus
  pending: string
  threshold: Threshold
}

export type Conditions = {
  [key in ConditionStatus]?: Condition
}

export type CheckStatus = 'active' | 'inactive'

export interface CheckBase {
  readonly id?: string
  readonly created: string
  readonly updated: string
  readonly orgID: string
  name: string
  desc: string
  query: string
  cron: string
  conditions: Condition[]
  status: CheckStatus

  // run state
  lastRunError: string
  lastRunStatus: string
  latestCompleted: string
  latestFailure: string
  latestScheduled: string
  latestSuccess: string
}

export interface Check extends Omit<CheckBase, 'conditions' | 'status'> {
  conditions: Conditions
  activeStatus: CheckStatus
  status: RemoteDataState
}

export interface CheckSortParams {
  direction: Sort
  type: SortTypes
  key: string
}
