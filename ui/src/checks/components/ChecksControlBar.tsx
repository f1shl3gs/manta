// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  PageControlBar,
  PageControlBarLeft,
  PageControlBarRight,
} from '@influxdata/clockface'
import SearchWidget from 'src/shared/components/SearchWidget'
import CreateCheckButton from 'src/checks/components/CreateCheckButton'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

// Actions
import {
  setCheckSearchTerm,
  setCheckSortOptions,
} from 'src/checks/actions/creators'

const mstp = (state: AppState) => {
  const {searchTerm, sortOptions} = state.resources[ResourceType.Checks]

  return {
    searchTerm,
    sortOptions,
  }
}

const mdtp = {
  setSearchTerm: setCheckSearchTerm,
  setSortOptions: setCheckSortOptions,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const ChecksControlBar: FunctionComponent<Props> = ({
  searchTerm,
  setSearchTerm,
  sortOptions,
  setSortOptions,
}) => {
  return (
    <PageControlBar fullWidth={false}>
      <PageControlBarLeft>
        <SearchWidget
          search={searchTerm}
          placeholder={'Filter checks...'}
          onSearch={setSearchTerm}
        />

        <ResourceSortDropdown
          resource={ResourceType.Checks}
          sortKey={sortOptions.key}
          sortType={sortOptions.type}
          sortDirection={sortOptions.direction}
          onSelect={setSortOptions}
        />
      </PageControlBarLeft>

      <PageControlBarRight>
        <CreateCheckButton />
      </PageControlBarRight>
    </PageControlBar>
  )
}

export default connector(ChecksControlBar)
