// Libraries
import React, {MouseEvent} from 'react'

// Components
import {Dropdown, Sort} from '@influxdata/clockface'

// Types
import {SortKey, SortTypes} from '../../types/Sort'

export interface SortDropdownItem {
  label: string
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

interface Props {
  width?: number

  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
  onSelect: (sortKey: SortKey, sortDirection: Sort, sortType: SortTypes) => void
}

const sortDropdownItems = [
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

const ResourceSortDropdown: React.FC<Props> = props => {
  const {sortKey, sortType, sortDirection, onSelect, width = 210} = props

  const {label} = sortDropdownItems.find(
    item =>
      item.sortKey === sortKey &&
      item.sortDirection === sortDirection &&
      item.sortType === sortType
  )!

  const button = (
    active: boolean,
    onClick: (e: MouseEvent<HTMLElement>) => void
  ) => (
    <Dropdown.Button active={active} onClick={onClick}>
      {`Sort by ${label}`}
    </Dropdown.Button>
  )

  const onItemClick = (item: SortDropdownItem): void => {
    const {sortKey, sortDirection, sortType} = item
    onSelect(sortKey, sortDirection, sortType)
  }

  const menu = (onCollapse?: () => void) => (
    <Dropdown.Menu onCollapse={onCollapse}>
      {sortDropdownItems.map(item => (
        <Dropdown.Item
          key={`${item.sortKey}${item.sortDirection}`}
          value={item}
          onClick={onItemClick}
          selected={
            item.sortKey === sortKey &&
            item.sortType === sortType &&
            item.sortDirection === sortDirection
          }
        >
          {item.label}
        </Dropdown.Item>
      ))}
    </Dropdown.Menu>
  )

  return (
    <Dropdown
      button={button}
      menu={menu}
      style={{flexBasis: `${width}px`, width: `${width}px`}}
    />
  )
}

export default ResourceSortDropdown
