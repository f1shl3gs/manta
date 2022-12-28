// Libraries
import React, {FunctionComponent, MouseEvent} from 'react'

// Components
import {Dropdown, Sort} from '@influxdata/clockface'

// Types
import {SortTypes} from 'src/types/sort'
import {ResourceType} from 'src/types/resources'

// Helpers
import {generateSortItems} from 'src/shared/components/resource_sort_dropdown/sortItems'

export interface SortDropdownItem {
  label: string
  sortKey: string
  sortType: SortTypes
  sortDirection: Sort
}

interface Props {
  width?: number

  resource: ResourceType
  sortKey: string
  sortType: SortTypes
  sortDirection: Sort
  onSelect: (sortKey: string, sortDirection: Sort, sortType: SortTypes) => void
}

const ResourceSortDropdown: FunctionComponent<Props> = ({
  sortKey,
  sortType,
  sortDirection,
  onSelect,
  resource,
  width = 210,
}) => {
  const sortDropdownItems = generateSortItems(resource)
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
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
