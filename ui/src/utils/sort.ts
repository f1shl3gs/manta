import {SortTypes} from 'types/sort'
import {get} from './object'
import {Sort} from '@influxdata/clockface'

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

export function getSortedResources<T>(
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