// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  Page,
  PageContents,
  PageControlBar,
  PageControlBarLeft,
  PageControlBarRight,
  PageHeader,
  PageTitle,
} from '@influxdata/clockface'
import CreateCheckButton from 'src/checks/components/CreateCheckButton'
import GetResources from 'src/resources/components/GetResources'
import SearchWidget from 'src/shared/components/SearchWidget'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

// Actions
import {setCheckSearchTerm} from 'src/checks/actions/creators'
import CheckCards from './components/CheckCards'

const mstp = (state: AppState) => {
  const {searchTerm, sortOptions} = state.resources[ResourceType.Checks]

  return {
    searchTerm,
    sortOptions,
  }
}

const mdtp = {
  setSearchTerm: setCheckSearchTerm,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const Inner: FunctionComponent<Props> = ({
  searchTerm,
  setSearchTerm,
  sortOptions,
}) => {
  return (
    <Page titleTag={'Checks'}>
      <PageHeader fullWidth={false}>
        <PageTitle title={'Checks'} />
      </PageHeader>

      <PageControlBar fullWidth={false}>
        <PageControlBarLeft>
          <SearchWidget
            search={searchTerm}
            placeholder={'Filter checks...'}
            onSearch={setSearchTerm}
          />
        </PageControlBarLeft>

        <PageControlBarRight>
          <CreateCheckButton />
        </PageControlBarRight>
      </PageControlBar>

      <PageContents>
        <CheckCards searchTerm={searchTerm} sortOption={sortOptions} />
      </PageContents>
    </Page>
  )
}

const ChecksIndex = connector(Inner)

export default () => {
  return (
    <GetResources resources={[ResourceType.Checks]}>
      <ChecksIndex />
    </GetResources>
  )
}
