import {Sort} from '@influxdata/clockface'
import {SortTypes} from 'src/types/sort'
import {ResourceType} from 'src/types/resources'

export const generateSortItems = (resourceType: ResourceType) => {
  switch (resourceType) {
    case ResourceType.Secrets:
      return [
        {
          label: 'Key (A → Z)',
          sortKey: 'key',
          sortType: SortTypes.String,
          sortDirection: Sort.Ascending,
        },
        {
          label: 'Key (Z → A)',
          sortKey: 'key',
          sortType: SortTypes.String,
          sortDirection: Sort.Descending,
        },
        {
          label: 'Modified (Oldest)',
          sortKey: 'updated',
          sortType: SortTypes.Date,
          sortDirection: Sort.Ascending,
        },
        {
          label: 'Modified (Newest)',
          sortKey: 'updated',
          sortType: SortTypes.Date,
          sortDirection: Sort.Descending,
        },
      ]

    default:
      return [
        {
          label: 'Name (A → Z)',
          sortKey: 'name',
          sortType: SortTypes.String,
          sortDirection: Sort.Ascending,
        },
        {
          label: 'Name (Z → A)',
          sortKey: 'name',
          sortType: SortTypes.String,
          sortDirection: Sort.Descending,
        },
        {
          label: 'Modified (Oldest)',
          sortKey: 'updated',
          sortType: SortTypes.Date,
          sortDirection: Sort.Ascending,
        },
        {
          label: 'Modified (Newest)',
          sortKey: 'updated',
          sortType: SortTypes.Date,
          sortDirection: Sort.Descending,
        },
      ]
  }
}
