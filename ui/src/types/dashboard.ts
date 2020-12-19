export interface Axis {
  bounds?: string[]
  label?: string
  prefix?: string
  suffix?: string
}

export interface Axes {
  x: Axis
  y: Axis
}

export interface XYViewProperties {
  type: "xy"
  timeFormat?: string
  axes: Axes
  shadeBelow?: boolean
}

export interface GaugeProperties {
  type: "gauge"
  prefix?: string
  suffix?: string
}

export type ViewProperties = XYViewProperties | GaugeProperties

export interface Cell {
  id: string

  name: string
  desc: string

  w?: number
  h?: number
  x?: number
  y?: number

  properties: ViewProperties
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

export type Dashboards = Dashboard[]