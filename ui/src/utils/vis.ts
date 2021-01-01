import { TimeZone } from './timeZone';
import { binaryPrefixFormatter, siPrefixFormatter, timeFormatter } from '@influxdata/giraffe';
import { DEFAULT_TIME_FORMAT, FORMAT_OPTIONS, resolveTimeFormat } from '../constants/timeFormat';

export type ColumnType = 'number' | 'string' | 'time' | 'boolean'

export type AxisScale = 'log' | 'linear'

export const VIS_SIG_DIGITS = 4

export interface Axis {
  bounds?: string[]
  label?: string
  prefix?: string
  suffix?: string
  base?: '' | '2' | '10'
  scale?: AxisScale
}

export type Base = Axis['base']

interface GetFormatterOptions {
  prefix?: string
  suffix?: string
  base?: Base
  timeZone?: TimeZone
  trimZeros?: boolean
  timeFormat?: string
}

export const getFormatter = (
  columnType: ColumnType,
  {
    prefix,
    suffix,
    base,
    timeZone,
    trimZeros = true,
    timeFormat = DEFAULT_TIME_FORMAT,
  }: GetFormatterOptions = {}
): null | ((x: any) => string) => {
  if (columnType === 'number' && base === '2') {
    return binaryPrefixFormatter({
      prefix,
      suffix,
      significantDigits: VIS_SIG_DIGITS,
    })
  }

  if (columnType === 'number') {
    return siPrefixFormatter({
      prefix,
      suffix,
      significantDigits: VIS_SIG_DIGITS,
      trimZeros,
    })
  }

  if (columnType === 'time') {
    return timeFormatter({
      timeZone: timeZone === 'Local' ? undefined : timeZone,
      format: resolveTimeFormat(timeFormat),
    })
  }

  return null
}