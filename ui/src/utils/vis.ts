import { TimeZone } from './timeZone';
import { binaryPrefixFormatter, siPrefixFormatter, timeFormatter } from '@influxdata/giraffe';

export type ColumnType = 'number' | 'string' | 'time' | 'boolean'

export const DEFAULT_TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss ZZ'

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

export const FORMAT_OPTIONS: Array<{text: string}> = [
  {text: DEFAULT_TIME_FORMAT},
  {text: 'DD/MM/YYYY HH:mm:ss.sss'},
  {text: 'MM/DD/YYYY HH:mm:ss.sss'},
  {text: 'YYYY/MM/DD HH:mm:ss'},
  {text: 'hh:mm a'},
  {text: 'HH:mm'},
  {text: 'HH:mm:ss'},
  {text: 'HH:mm:ss ZZ'},
  {text: 'HH:mm:ss.sss'},
  {text: 'MMMM D, YYYY HH:mm:ss'},
  {text: 'dddd, MMMM D, YYYY HH:mm:ss'},
]

export const resolveTimeFormat = (timeFormat: string) => {
  if (FORMAT_OPTIONS.find(d => d.text === timeFormat)) {
    return timeFormat
  }

  return DEFAULT_TIME_FORMAT
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