// Libraries
import React from 'react'

// Components
import EmptyEndpointsList from './EmptyEndpointsList'
import NotificationEndpointCard from './NotificationEndpointCard'
import {ResourceList, Sort} from '@influxdata/clockface'

// Types
import {NotificationEndpoint} from 'client'
import {SortKey, SortTypes} from 'types/sort'

// Utils
import {getSortedResources} from 'utils/sort'

interface Props {
  search: string
  endpoints: NotificationEndpoint[]
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

const NotificationEndpointCards: React.FC<Props> = props => {
  const {endpoints, search, sortKey, sortType, sortDirection} = props

  const body = (filtered: NotificationEndpoint[]) => (
    <ResourceList.Body emptyState={<EmptyEndpointsList search={search} />}>
      {getSortedResources<NotificationEndpoint>(
        filtered,
        sortKey,
        sortType,
        sortDirection
      ).map(endpoint => (
        <NotificationEndpointCard key={endpoint.id} endpoint={endpoint} />
      ))}
    </ResourceList.Body>
  )

  return <ResourceList>{body(endpoints)}</ResourceList>
}

export default NotificationEndpointCards
