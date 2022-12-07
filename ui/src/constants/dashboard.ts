import {Sort} from '@influxdata/clockface'
import {DashboardSortKey, SortTypes} from 'src/types/sort'

export const MIN_DECIMAL_PLACES = 0
export const MAX_DECIMAL_PLACES = 10

export const DEFAULT_DASHBOARD_SORT_OPTIONS = {
  direction: Sort.Ascending,
  type: SortTypes.String,
  key: 'name' as DashboardSortKey,
}
