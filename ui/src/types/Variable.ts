import {Common} from './Common'

export interface Variable extends Common {
  name: string
  desc?: string
  type: 'query' | 'static'
  value: string
}
