// Libraries
import {RemoteDataState} from '@influxdata/clockface'

// Types
import {
  Axis,
  Base,
  GaugeViewProperties,
  LinePlusSingleStatViewProperties,
  SingleStatViewProperties,
  ViewProperties,
  ViewType,
  XYViewProperties,
} from 'src/types/cells'
import {DashboardQuery} from 'src/types/dashboards'
import {Color} from 'src/types/colors'

// Constants
import {DEFAULT_THRESHOLDS_LIST_COLORS} from 'src/constants/thresholds'

export const defaultView = () => {
  return {
    name,
    note: '',
    status: RemoteDataState.Done,

    colors: new Array<Color>(),
    queries: new Array<DashboardQuery>(),
  }
}

const defaultSingleStatViewProperties = () => ({
  showNoteWhenEmpty: false,
  colors: DEFAULT_THRESHOLDS_LIST_COLORS as Color[],
  prefix: '',
  suffix: '',
  tickPrefix: '',
  tickSuffix: '',
  decimalPlaces: {
    isEnforced: true,
    digits: 2,
  },
})

const defaultLineViewProperties = () => ({
  geom: 'line',
  xColumn: null,
  yColumn: null,
  position: 'overlaid',
  axes: {
    x: {
      bounds: ['', ''],
      label: '',
      prefix: '',
      suffix: '',
      base: '10',
      scale: 'linear',
    } as Axis,
    y: {
      bounds: ['', ''],
      label: '',
      prefix: '',
      suffix: '',
      base: '10' as Base,
      scale: 'linear',
    } as Axis,
  },
})

const defaultGaugeViewProperties = () => ({
  prefix: '',
  tickPrefix: '',
  suffix: '',
  tickSuffix: '',
  showNoteWhenEmpty: false,
  decimalPlaces: {
    isEnforced: true,
    digits: 2,
  },
})

const NEW_VIEW_CREATORS = {
  xy: (): XYViewProperties => ({
    ...defaultView(),
    ...defaultLineViewProperties(),

    type: 'xy',
    position: 'overlaid',
    geom: 'line',
  }),

  'line-plus-single-stat': (): LinePlusSingleStatViewProperties => ({
    ...defaultView(),
    ...defaultLineViewProperties(),
    ...defaultSingleStatViewProperties(),

    type: 'line-plus-single-stat',
    position: 'overlaid',
  }),

  'single-stat': (): SingleStatViewProperties => ({
    ...defaultView(),
    ...defaultSingleStatViewProperties(),

    type: 'single-stat',
  }),

  gauge: (): GaugeViewProperties => ({
    ...defaultView(),
    ...defaultGaugeViewProperties(),
    type: 'gauge',
  }),
}

export function createView(viewType: ViewType): ViewProperties {
  const creator = NEW_VIEW_CREATORS[viewType]

  if (!creator) {
    throw new Error(`no view properties creator implemented for ${viewType}`)
  }

  return creator()
}
