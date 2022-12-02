export interface Variable {
  id: string
  created: string
  updated: string
  orgID: string
  name: string
  desc?: string
  type: 'query' | 'static'
  value: string
}
