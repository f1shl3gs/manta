import constate from 'constate'
import {Sort} from '@influxdata/clockface'
import {SortKey} from '../../types/Sort'

interface Sortable {
  name: string
  created: string
  updated: string
}

const sortBy = <T extends Sortable>(
  resources: T[],
  sort: Sort,
  sortKey: SortKey,
  sortDirection: Sort
): T[] => {
  return resources
}
