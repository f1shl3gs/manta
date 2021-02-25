import {Common} from './Common'

export interface Condition {}

export interface Check extends Common {
  name: string
  desc?: string
  status: string
  conditions: Condition[]

  // status
  latestCompleted: string
  latestScheduled: string
  latestSuccess: string
  latestFailure: string
  lastRunStatus: string
  lastRunError?: string
}
