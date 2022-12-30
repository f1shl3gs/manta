// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'
import TabbedPageHeader from 'src/shared/components/TabbedPageHeader'
import CreateNotificationEndpointButton from 'src/notification_endpoints/components/CreateNotificationEndpointButton'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'

// Actions
import {
  setNotificationEndpointSearchTerm,
  setNotificationEndpointSortOptions,
} from 'src/notification_endpoints/actions/creators'

const mstp = (state: AppState) => {
  const {searchTerm, sortOptions} =
    state.resources[ResourceType.NotificationEndpoints]

  return {
    searchTerm,
    sortOptions,
  }
}

const mdtp = {
  setSearchTerm: setNotificationEndpointSearchTerm,
  setSortOptions: setNotificationEndpointSortOptions,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const ContentHeader: FunctionComponent<Props> = ({
  searchTerm,
  setSearchTerm,
  sortOptions,
  setSortOptions,
}) => {
  const left = (
    <>
      <SearchWidget
        search={searchTerm}
        placeholder={'Filter notification endpoints...'}
        onSearch={setSearchTerm}
      />

      <ResourceSortDropdown
        resource={ResourceType.NotificationEndpoints}
        sortKey={sortOptions.key}
        sortType={sortOptions.type}
        sortDirection={sortOptions.direction}
        onSelect={setSortOptions}
      />
    </>
  )

  return (
    <TabbedPageHeader
      left={left}
      right={<CreateNotificationEndpointButton />}
    />
  )
}

export default connector(ContentHeader)
