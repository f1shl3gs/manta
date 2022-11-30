import {Sort} from '@influxdata/clockface'

import {SortTypes} from 'src/types/Sort'
import {get} from 'lodash'

function sortBy<T>(
  resourceList: T[],
  iter: (t: T) => string | number,
  direction: Sort
): T[] {
  if (resourceList == null) {
    return []
  }

  const sorted = resourceList.sort((a, b) => {
    const av = iter(a)
    const bv = iter(b)

    return av > bv ? 1 : bv > av ? -1 : 0
  })

  if (direction === Sort.Descending) {
    return sorted.reverse()
  }

  return sorted
}

function orderByType(data: string, type: SortTypes) {
  switch (type) {
    case SortTypes.String:
      return data.toLowerCase()
    case SortTypes.Date:
      return Date.parse(data)
    case SortTypes.Float:
      return parseFloat(data)
    default:
      return data
  }
}

interface Sortable {
  name: string
  updated: string
}

export function getSortedResources<T extends Sortable>(
  resourceList: T[],
  sortKey: string,
  sortType: SortTypes,
  sortDirection: Sort
): T[] {
  if (sortKey && sortDirection) {
    return sortBy<T>(
      resourceList,
      r => orderByType(get(r, sortKey), sortType),
      sortDirection
    )
  }

  return resourceList
}
