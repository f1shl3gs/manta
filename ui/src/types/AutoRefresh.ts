export enum AutoRefreshOptionType {
  Option = 'option',
  Header = 'header',
}

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
  id: string
  label: string
  type: AutoRefreshOptionType,
  seconds: number
}
