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
import GetResources from 'src/resources/components/GetResources'
import {ResourceType} from 'src/types/resources'
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import CreateDashboardButton from 'src/dashboards/components/CreateDashboardButton'
import DashboardCards from 'src/dashboards/components/DashboardCards'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {
  setDashboardSearchTerm,
  setDashboardSort,
} from 'src/dashboards/actions/creators'

type ReduxProps = ConnectedProps<typeof connector>
type Props = ReduxProps

export const DashboardsIndex: FunctionComponent<Props> = ({
  searchTerm,
  sortOptions,
  setDashboardSearchTerm,
  setDashboardSort,
}) => {
  return (
    <>
      <Page titleTag="Dashboards">
        <PageHeader fullWidth={false}>
          <PageTitle title="Dashboards" />
        </PageHeader>

        <PageControlBar fullWidth={false}>
          <PageControlBarLeft>
            <SearchWidget
              onSearch={setDashboardSearchTerm}
              placeholder="Filter dashboards..."
              search={searchTerm}
            />
            <ResourceSortDropdown
              sortKey={sortOptions.key}
              sortType={sortOptions.type}
              sortDirection={sortOptions.direction}
              onSelect={(key, direction, type) => {
                setDashboardSort({
                  key,
                  type,
                  direction,
                })
              }}
            />
          </PageControlBarLeft>

          <PageControlBarRight>
            <CreateDashboardButton />
          </PageControlBarRight>
        </PageControlBar>

        <PageContents>
          <DashboardCards search={searchTerm} sortOption={sortOptions} />
        </PageContents>
      </Page>
    </>
  )
}

const mstp = (state: AppState) => {
  const {sortOptions, searchTerm} = state.resources[ResourceType.Dashboards]

  return {
    sortOptions,
    searchTerm,
  }
}

const mdtp = {
  setDashboardSearchTerm,
  setDashboardSort,
}

const connector = connect(mstp, mdtp)

const ToExport = connector(DashboardsIndex)

export default connector(() => (
  <GetResources resources={[ResourceType.Dashboards]}>
    <ToExport />
  </GetResources>
))
