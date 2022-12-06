// Libraries
import React, {FunctionComponent, lazy} from 'react'
import {Route, Routes } from 'react-router-dom'
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
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import CreateDashboardButton from 'src/dashboards/components/CreateDashboardButton'
import DashboardCards from 'src/dashboards/components/DashboardCards'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'

// Actions
import {
  setDashboardSearchTerm,
  setDashboardSort,
} from 'src/dashboards/actions/creators'

// Lazy load
const DashboardImportOverlay = lazy(() => import('src/dashboards/components/DashboardImportOverlay'))

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

// /orgs/:orgID/dashboards
export default () => (
  <GetResources resources={[ResourceType.Dashboards]}>
    <Routes>
      <Route index element={<ToExport />} />
      <Route
        path="import"
        element={<DashboardImportOverlay />}
      />
    </Routes>
  </GetResources>
)
