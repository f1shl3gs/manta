import {Sort} from '@influxdata/clockface'
import {ViewProperties} from 'src/types/cells'
import {DashboardSortKey, SortTypes} from 'src/types/sort'

export const MIN_DECIMAL_PLACES = 0
export const MAX_DECIMAL_PLACES = 10

export const DEFAULT_VIEWPROPERTIES: ViewProperties = {
  type: 'xy',
  xColumn: '_time',
  yColumn: '_value',
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

export const DEFAULT_DASHBOARD_SORT_OPTIONS = {
  direction: Sort.Ascending,
  type: SortTypes.String,
  key: 'name' as DashboardSortKey,
}
