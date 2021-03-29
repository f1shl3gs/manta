export interface Variable {
  id: string
  created: string
  updated: string
  name: string
  desc: string
  orgID: string
  type: 'query' | 'static'
  value: string
}
