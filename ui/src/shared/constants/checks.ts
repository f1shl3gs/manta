// Libraries
import {Sort} from '@influxdata/clockface'

// Types
import {CheckSortKey} from 'src/types/checks'
import {SortTypes} from 'src/types/sort'

export const DEFAULT_CHECK_SORT_OPTIONS = {
  direction: Sort.Ascending,
  type: SortTypes.String,
  key: 'name' as CheckSortKey,
}
