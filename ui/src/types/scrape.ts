export interface Scrape {
  id?: string
  created?: string
  updated: string
  name: string
  desc?: string
  orgID: string

  labels?: {[key: string]: string}
  targets: string[]
}
