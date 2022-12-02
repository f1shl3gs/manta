// Libraries
import React, {FC, useState} from 'react'

// Components
import {
  Page,
  PageContents,
  PageControlBar,
  PageControlBarLeft,
  PageControlBarRight,
  PageHeader,
  PageTitle,
  Sort,
} from '@influxdata/clockface'
import {GetResources, ResourceType} from 'src/shared/components/GetResources'
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import {SortKey, SortTypes} from 'src/types/sort'
import CreateDashboardButton from 'src/dashboards/CreateDashboardButton'
import DashboardCards from 'src/dashboards/DashboardCards'

export const DashboardsPage: FC = () => {
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  return (
    <>
      <Page titleTag="Dashboards">
        <PageHeader fullWidth={false}>
          <PageTitle title="Dashboards" />
        </PageHeader>

        <PageControlBar fullWidth={false}>
          <PageControlBarLeft>
            <SearchWidget
              onSearch={t => setSearch(t)}
              placeholder="Filter dashboards..."
              search={search}
            />
            <ResourceSortDropdown
              sortKey={sortOption.key}
              sortType={sortOption.type}
              sortDirection={sortOption.direction}
              onSelect={(key, direction, type) => {
                setSortOption({
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
          <GetResources type={ResourceType.Dashboards}>
            <DashboardCards search={search} sortOption={sortOption} />
          </GetResources>
        </PageContents>
      </Page>
    </>
  )
}

export default DashboardsPage
