// Libraries
import React from 'react'

// Components
import {ResourceList, Sort} from '@influxdata/clockface'
import EmptyChecksList from './EmptyChecksList'
import CheckCard from './CheckCard'

// Types
import {Check} from 'types/Check'
import {SortKey, SortTypes} from 'types/sort'

// Utils
import {getSortedResources} from 'utils/sort'

interface Props {
  search: string
  checks: Check[]
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

const CheckCards: React.FC<Props> = props => {
  const {checks, search, sortKey, sortType, sortDirection} = props

  const body = (filtered: Check[]) => (
    <ResourceList.Body emptyState={<EmptyChecksList search={search} />}>
      {getSortedResources<Check>(
        filtered,
        sortKey,
        sortType,
        sortDirection
      ).map(check => (
        <CheckCard key={check.id} check={check} />
      ))}
    </ResourceList.Body>
  )

  return <ResourceList>{body(checks)}</ResourceList>
}

export default CheckCards
