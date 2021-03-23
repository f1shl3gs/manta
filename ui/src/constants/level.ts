import {ComponentColor, InfluxColors} from '@influxdata/clockface'

export const LEVEL_COLORS = {
  OK: InfluxColors.Viridian,
  INFO: InfluxColors.Ocean,
  WARN: InfluxColors.Thunder,
  CRIT: InfluxColors.Fire,
  UNKNOWN: InfluxColors.Amethyst,
}

export const LEVEL_COMPONENT_COLORS = {
  OK: ComponentColor.Success,
  INFO: ComponentColor.Primary,
  WARN: ComponentColor.Warning,
  CRIT: ComponentColor.Danger,
  UNKNOWN: ComponentColor.Secondary,
}
