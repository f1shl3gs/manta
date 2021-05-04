// Libraries
import React from 'react'

// Components
import {Otcl} from '../../types/otcl'
import {ResourceList, Sort} from '@influxdata/clockface'
import EmptyOtclList from './EmptyOtclList'
import OtclCard from './OtclCard'

// Types
import {SortKey, SortTypes} from '../../types/sort'

// Utils
import {getSortedResources} from '../../utils/sort'

interface Props {
  list: Otcl[]
  onDelete: (otcl: Otcl) => void

  searchTerm: string
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

const OtclCards: React.FC<Props> = props => {
  const {list, searchTerm, sortKey, sortType, sortDirection, onDelete} = props

  const body = (filtered: Otcl[]) => (
    <ResourceList.Body emptyState={<EmptyOtclList searchTerm={searchTerm} />}>
      {getSortedResources<Otcl>(filtered, sortKey, sortType, sortDirection).map(
        otcl => (
          <OtclCard key={otcl.id} otcl={otcl} onDelete={onDelete} />
        )
      )}
    </ResourceList.Body>
  )

  return <ResourceList>{body(list)}</ResourceList>
}

export default OtclCards
