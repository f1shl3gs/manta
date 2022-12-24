// Types
import {ComponentColor, InfluxColors} from '@influxdata/clockface'

export const DEFAULT_CHECK_NAME = 'Name this check'
export const DEFAULT_CHECK_DESC = 'Describe this check'

export const DEFAULT_CHECK_CRON = '@every 1m'

export const LEVEL_COMPONENT_COLORS = {
  ok: ComponentColor.Success,
  info: ComponentColor.Primary,
  warn: ComponentColor.Warning,
  crit: ComponentColor.Danger,
}

export const CHECK_STATUS_COLORS = {
  ok: InfluxColors.Viridian,
  info: InfluxColors.Ocean,
  warn: InfluxColors.Thunder,
  crit: InfluxColors.Fire,
  unknown: InfluxColors.Amethyst,
}
