// Libraries
import React, {FunctionComponent} from 'react'
import {ConnectedProps, connect} from 'react-redux'

// Components
import {
  PageControlBar,
  PageControlBarLeft,
  PageControlBarRight,
} from '@influxdata/clockface'
import CreateDashboardButton from 'src/dashboards/components/CreateDashboardButton'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'
import SearchWidget from 'src/shared/components/SearchWidget'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

// Actions
import {
  setDashboardSearchTerm,
  setDashboardSortOptions,
} from 'src/dashboards/actions/creators'

const mstp = (state: AppState) => {
  const {sortOptions, searchTerm} = state.resources[ResourceType.Dashboards]

  return {
    sortOptions,
    searchTerm,
  }
}

const mdtp = {
  setSearchTerm: setDashboardSearchTerm,
  setSortOptions: setDashboardSortOptions,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const DashboardsPageHeader: FunctionComponent<Props> = ({
  sortOptions,
  searchTerm,
  setSearchTerm,
  setSortOptions,
}) => {
  return (
    <PageControlBar fullWidth={false}>
      <PageControlBarLeft>
        <SearchWidget
          onSearch={setSearchTerm}
          placeholder="Filter dashboards..."
          search={searchTerm}
        />
        <ResourceSortDropdown
          resource={ResourceType.Dashboards}
          sortKey={sortOptions.key}
          sortType={sortOptions.type}
          sortDirection={sortOptions.direction}
          onSelect={setSortOptions}
        />
      </PageControlBarLeft>

      <PageControlBarRight>
        <CreateDashboardButton />
      </PageControlBarRight>
    </PageControlBar>
  )
}

export default connector(DashboardsPageHeader)
