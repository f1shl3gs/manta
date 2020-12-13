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

export interface Panel {
  name: string
  desc: string

  w?: number
  h?: number
  x?: number
  y?: number

  properties?: ViewProperties
}

export type Panels = Panel[]

export interface Dashboard {
  id: string
  created: string
  updated: string
  name: string
  desc: string
  orgID: string
  panels: Panels
}

