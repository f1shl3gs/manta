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
import {Color} from 'src/types/colors'

// Constants
import {DEFAULT_THRESHOLDS_LIST_COLORS} from 'src/constants/thresholds'
import {LineHoverDimension} from '@influxdata/giraffe/dist/types'
import {DEFAULT_LINE_COLORS} from 'src/constants/graphColorPalettes'
import {DEFAULT_GAUGE_COLORS} from '@influxdata/giraffe'

const tickProps = {
  generateXAxisTicks: [],
  generateYAxisTicks: [],
  xTotalTicks: null,
  xTickStart: null,
  xTickStep: null,
  yTotalTicks: null,
  yTickStart: null,
  yTickStep: null,
}

export const defaultView = () => {
  return {
    name,
    note: '',
    status: RemoteDataState.Done,

    colors: new Array<Color>(),
    queries: [
      {
        name: 'query 1',
        text: '',
        hidden: false,
      },
    ],
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
  xColumn: '_time',
  yColumn: '_value',
  colors: DEFAULT_LINE_COLORS as Color[],
  showNoteWhenEmpty: false,
  ...tickProps,
  hoverDimension: 'auto' as LineHoverDimension,
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
  colors: DEFAULT_GAUGE_COLORS as Color[],
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
    geom: 'line',
    position: 'overlaid',
    xColumn: null,
    yColumn: null,
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
