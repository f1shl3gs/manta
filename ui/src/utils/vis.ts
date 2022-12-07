import {TimeZone} from 'src/utils/timeZone'
import {
  binaryPrefixFormatter,
  siPrefixFormatter,
  timeFormatter,
} from '@influxdata/giraffe'
import {DEFAULT_TIME_FORMAT, resolveTimeFormat} from 'src/constants/timeFormat'

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
    timeZone = 'Local',
    trimZeros = true,
    timeFormat = DEFAULT_TIME_FORMAT,
  }: GetFormatterOptions = {}
): null | ((x: any) => string) => {
  if (columnType === 'number' && base === '2') {
    return binaryPrefixFormatter({
      prefix,
      suffix,
      significantDigits: VIS_SIG_DIGITS,
      format: true,
    })
  }

  if (columnType === 'number' && base === '10') {
    return siPrefixFormatter({
      prefix,
      suffix,
      significantDigits: VIS_SIG_DIGITS,
      trimZeros,
      format: true,
    })
  }

  if (columnType === 'number' && base === '') {
    return siPrefixFormatter({
      prefix,
      suffix,
      significantDigits: VIS_SIG_DIGITS,
      trimZeros,
      format: false,
    })
  }

  if (columnType === 'time') {
    const formatOptions = {
      timeZone: timeZone === 'Local' ? undefined : timeZone,
      format: resolveTimeFormat(timeFormat),
    }
    if (formatOptions.format.includes('HH')) {
      formatOptions['hour12'] = false
    }
    return timeFormatter(formatOptions)
  }

  return null
}
