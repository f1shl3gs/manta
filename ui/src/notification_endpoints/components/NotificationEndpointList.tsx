// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import EmptyResources from 'src/resources/components/EmptyResources'
import CreateNotificationEndpointButton from 'src/notification_endpoints/components/CreateNotificationEndpointButton'
import NotificationEndpointCard from 'src/notification_endpoints/components/NotificationEndpointCard'
import FilterList from 'src/shared/components/FilterList'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'

// Utils
import {getSortedResources} from 'src/shared/utils/sort'

// Selectors
import {getAll} from 'src/resources/selectors'

const NotificationEndpointList: FunctionComponent = () => {
  const {endpoints, searchTerm, sortOptions} = useSelector(
    (state: AppState) => {
      const endpoints = getAll<NotificationEndpoint>(
        state,
        ResourceType.NotificationEndpoints
      )
      const {searchTerm, sortOptions} =
        state.resources[ResourceType.NotificationEndpoints]

      return {
        endpoints,
        searchTerm,
        sortOptions,
      }
    }
  )

  return (
    <FilterList<NotificationEndpoint>
      list={endpoints}
      search={searchTerm}
      searchKeys={['name', 'desc', 'url', 'contentTemplate']}
    >
      {filtered => {
        if (filtered && filtered.length === 0) {
          return (
            <EmptyResources
              resource={ResourceType.NotificationEndpoints}
              createButton={<CreateNotificationEndpointButton />}
            />
          )
        }

        return (
          <div>
            {getSortedResources(
              filtered,
              sortOptions.key,
              sortOptions.type,
              sortOptions.direction
            ).map(ep => (
              <NotificationEndpointCard key={ep.id} notificationEndpoint={ep} />
            ))}
          </div>
        )
      }}
    </FilterList>
  )
}

export default NotificationEndpointList
