import {Sort} from '@influxdata/clockface'
import {SecretSortKey} from 'src/types/secrets'
import {SortTypes} from 'src/types/sort'

export const DEFAULT_SECRET_SORT_OPTIONS = {
  direction: Sort.Ascending,
  type: SortTypes.String,
  key: 'key' as SecretSortKey,
}
