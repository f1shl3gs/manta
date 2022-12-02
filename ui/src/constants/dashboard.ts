import {ViewProperties} from 'src/types/dashboard'

export const MIN_DECIMAL_PLACES = 0
export const MAX_DECIMAL_PLACES = 10

export const defaultViewProperties: ViewProperties = {
  type: 'xy',
  xColumn: 'time',
  yColumn: 'value',
  hoverDimension: 'auto',
  geom: 'line',
  position: 'overlaid',
  axes: {
    x: {},
    y: {},
  },
  queries: [
    {
      name: 'query 1',
      text: '',
      hidden: false,
    },
  ],
}
