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

export interface GaugeProperties {
  type: 'gauge'
  note?: string
  prefix: string
  suffix: string
  legend: Legend
  decimalPlaces: DecimalPlaces
  queries: DashboardQuery[]
}

export type ViewProperties = XYViewProperties | GaugeProperties

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

export type NewView<T extends ViewProperties = ViewProperties> = Omit<View<T>,
  'id'>

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