export enum AutoRefreshStatus {
  Active = 'active',
  Paused = 'paused',
  Disabled = 'disabled',
}

export interface AutoRefresh {
  status: AutoRefreshStatus
  interval: number
}

export interface AutoRefreshOption {
  label: string
  seconds: number
}
