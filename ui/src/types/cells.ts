import {RemoteDataState} from '@influxdata/clockface'
import {LineInterpolation} from '@influxdata/giraffe'
import {DashboardColor} from 'src/types/colors'
import {DashboardQuery} from 'src/types/dashboards'

export interface Axis {
  bounds?: string[]
  label?: string
  prefix?: string
  suffix?: string
  base?: '' | '2' | '10'
}

export interface Axes {
  x: Axis
  y: Axis
}

export type XYGeom = 'line' | 'step' | 'stacked' | 'bar' | 'monotoneX'

export interface XYViewProperties {
  type: 'xy'
  timeFormat?: string
  axes: Axes
  shadeBelow?: boolean
  xColumn?: string
  yColumn?: string
  hoverDimension?: 'auto' | 'x' | 'y' | 'xy'
  position: 'overlaid' | 'stacked'
  geom: XYGeom
  queries: DashboardQuery[]
  colors: DashboardColor[]
  interpolation: LineInterpolation
}

export interface Legend {
  type?: 'static'
  orientation?: 'top' | 'bottom' | 'left' | 'right'
}

export interface DecimalPlaces {
  isEnforced?: boolean
  digits?: number
}

export interface GaugeViewProperties {
  type: 'gauge'
  note?: string
  prefix: string
  suffix: string
  tickPrefix: string
  tickSuffix: string
  decimalPlaces: DecimalPlaces
  queries: DashboardQuery[]
  colors: DashboardColor[]
}

export interface SingleStatViewProperties {
  type: 'single-stat'
  queries: DashboardQuery[]
  colors: DashboardColor[]
  note: string
  showNoteWhenEmpty: boolean
  prefix: string
  tickPrefix: string
  suffix: string
  tickSuffix: string
  legend?: Legend
  decimalPlaces: DecimalPlaces
}

export interface HistogramViewProperties {
  type: 'histogram'
  queries: DashboardQuery[]
  note: string
  showNoteWhenEmpty: boolean
  xColumn: string
  fillColumns: string[]
  xDomain: number[]
  xAxisLabel: string
  position: 'overlaid' | 'stacked'
}

export interface MarkdownViewProperties {
  type: 'markdown'
  note: string

  // we don't need this actually
  queries: DashboardQuery[]
}

export interface BandViewProperties {
  type: 'band'
  timeFormat?: string
  queries: DashboardQuery[]
  note: string
  showNoteWhenEmpty: boolean
  axes: Axes
  legend: Legend
  xColumn?: string
  hoverDimension?: 'auto' | 'x' | 'y' | 'xy'
  geom: XYGeom
}

export interface TableViewProperties {
  type: 'table'
  queries: DashboardQuery[]
  note: string
  showNoteWhenEmpty: boolean
  tableOptions: {
    verticalTimeAxis?: boolean
  }
  timeFormat: string
  decimalPlaces: DecimalPlaces
}

export interface ScatterViewProperties {
  type: 'scatter'
  timeFormat?: string
  queries: DashboardQuery[]
  note: string
  showNoteWhenEmpty: boolean
  xColumn: string
  xPrefix: string
  xSuffix: string
  yColumn: string
  yPrefix: string
  ySuffix: string
}

export interface LinePlusSingleStatViewProperties {
  type: 'line-plus-single-stat'
  queries: DashboardQuery[]

  colors: DashboardColor[]
  axes: Axes
  xColumn?: string
  generateXAxisTicks?: string[]
  xTotalTicks?: number
  xTickStart?: number
  xTickStep?: number
  yColumn?: string
  generateYAxisTicks?: string[]
  yTotalTicks?: number
  yTickStart?: number
  yTickStep?: number
  shadeBelow?: boolean
  hoverDimension?: 'auto' | 'x' | 'y' | 'xy'
  position: 'overlaid' | 'stacked'
  prefix: string
  suffix: string
  decimalPlaces: DecimalPlaces
  timeFormat?: string
}

export interface MosaicViewProperties {
  type: 'mosaic'
  queries: DashboardQuery[]
}

export interface HeatmapViewProperties {
  type: 'heatmap'
  queries: DashboardQuery[]
}

export type ViewProperties =
  | XYViewProperties
  | GaugeViewProperties
  | SingleStatViewProperties
  | HistogramViewProperties
  | MarkdownViewProperties
  | BandViewProperties
  | TableViewProperties
  | ScatterViewProperties
  | LinePlusSingleStatViewProperties
  | MosaicViewProperties
  | HeatmapViewProperties

export type ViewType = ViewProperties['type']

export interface Cell {
  id: string

  name: string
  desc: string

  w?: number
  h?: number
  x?: number
  y?: number

  minH?: number
  minW?: number
  maxW?: number

  viewProperties: ViewProperties

  // extra information
  dashboardID: string
  status: RemoteDataState
}

export type NewCell = Omit<Cell, 'id'>

export type Cells = Cell[]

export type Base = Axis['base']

export type AxisScale = 'log' | 'linear'

// CellEntities defines the result of normalizr's normalization of the
// 'cells' resource
export interface CellEntities {
  cells: {
    [uuid: string]: Cell
  }
}
