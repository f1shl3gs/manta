import {Common} from './Common'

export type CheckStatusLevel = 'UNKNOWN' | 'OK' | 'INFO' | 'WARN' | 'CRIT'

export interface Condition {
  status: CheckStatusLevel
  pending?: string
  threshold: Threshold
}

export interface CheckStatus {
  latestCompleted: string
  latestScheduled: string
  latestSuccess: string
  latestFailure: string
  lastRunStatus: string
  lastRunError?: string
}

export interface Check extends Common, CheckStatus {
  name: string
  desc?: string
  status: string
  expr: string
  conditions: Condition[]
}

export interface ThresholdBase {
  level?: string
}

export type GreatThanThreshold = ThresholdBase & {
  type: 'gt'
  value: number
}

export type LessThanThreshold = ThresholdBase & {
  type: 'lt'
  value: number
}

export type InsideThreshold = ThresholdBase & {
  type: 'inside'
  max: number
  min: number
}

export type OutsideThreshold = ThresholdBase & {
  type: 'outside'
  max: number
  min: number
}

export type Threshold =
  | GreatThanThreshold
  | LessThanThreshold
  | InsideThreshold
  | OutsideThreshold

export type ThresholdType = Threshold['type']
