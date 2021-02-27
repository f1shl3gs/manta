import {Common} from './Common'

export interface Condition {}

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
