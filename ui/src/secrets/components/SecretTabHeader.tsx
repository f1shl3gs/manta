// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Componetns
import CreateSecretButton from './CreateSecretButton'
import TabbedPageHeader from 'src/shared/components/TabbedPageHeader'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import SearchWidget from 'src/shared/components/SearchWidget'
import {setSecretSearchTerm, setSecretSortOptions} from '../actions/creators'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'

const mstp = (state: AppState) => {
  const {searchTerm, sortOptions} = state.resources[ResourceType.Secrets]

  return {
    searchTerm,
    sortOptions,
  }
}

const mdtp = {
  setSearchTerm: setSecretSearchTerm,
  setSortOptions: setSecretSortOptions,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const SecretTabHeader: FunctionComponent<Props> = ({
  searchTerm,
  setSearchTerm,
  sortOptions,
  setSortOptions,
}) => {
  const left = (
    <>
      <SearchWidget
        search={searchTerm}
        placeholder={'Filter secrets...'}
        onSearch={setSearchTerm}
      />

      <ResourceSortDropdown
        resource={ResourceType.Secrets}
        sortKey={sortOptions.key}
        sortType={sortOptions.type}
        sortDirection={sortOptions.direction}
        onSelect={setSortOptions}
      />
    </>
  )

  return <TabbedPageHeader left={left} right={<CreateSecretButton />} />
}

export default connector(SecretTabHeader)
