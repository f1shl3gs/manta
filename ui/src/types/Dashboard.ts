import {DashboardColor} from 'src/types/Colors'

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
  legend: Legend
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
  legend: Legend
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

export interface GenView {
  readonly id?: string
  name: string
  properties: ViewProperties
}

export interface View<T extends ViewProperties = ViewProperties>
  extends GenView {
  properties: T
  cellID?: string
  dashboardID?: string
}

export type NewView<T extends ViewProperties = ViewProperties> = Omit<
  View<T>,
  'id'
>

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

  viewProperties?: ViewProperties
}

export type Cells = Cell[]

export interface Dashboard {
  id: string
  created: string
  updated: string
  name: string
  desc: string
  orgID: string
  cells: Cells
}

export interface DashboardQuery {
  name?: string
  text: string
  hidden: boolean
}

export type Dashboards = Dashboard[]

export type Base = Axis['base']

export type AxisScale = 'log' | 'linear'
