import {Common} from './Common'

export interface Condition {}

export interface Check extends Common {
  name: string
  desc?: string
  status: string
  conditions: Condition[]
}
