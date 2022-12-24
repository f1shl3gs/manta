import {AutoRefreshOption, AutoRefreshOptionType} from 'src/types/autoRefresh'

export const DROPDOWN_WIDTH_COLLAPSED = 50
export const DROPDOWN_WIDTH_FULL = 100

export const AutoRefreshDropdownOptions: AutoRefreshOption[] = [
  {
    id: 'pause',
    type: AutoRefreshOptionType.Option,
    label: 'Pause',
    seconds: 0,
  },
  {
    id: '5s',
    type: AutoRefreshOptionType.Option,
    label: '5s',
    seconds: 5,
  },
  {
    id: '15s',
    type: AutoRefreshOptionType.Option,
    label: '15s',
    seconds: 15,
  },
  {
    id: '30s',
    type: AutoRefreshOptionType.Option,
    label: '30s',
    seconds: 30,
  },
  {
    id: '60s',
    type: AutoRefreshOptionType.Option,
    label: '60s',
    seconds: 60,
  },
  {
    id: '5m',
    type: AutoRefreshOptionType.Option,
    label: '5m',
    seconds: 5 * 60,
  },
  {
    id: '15m',
    type: AutoRefreshOptionType.Option,
    label: '15m',
    seconds: 15 * 60,
  },
]
