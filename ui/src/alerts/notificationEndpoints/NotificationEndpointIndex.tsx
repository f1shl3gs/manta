// Libraries
import React, {useState} from 'react'

// Components
import SearchWidget from '../../shared/components/SearchWidget'
import ResourceSortDropdown from '../../shared/components/ResourceSortDropdown'
import {
  Button,
  Columns,
  ComponentColor,
  ComponentStatus,
  Grid,
  IconFont,
  Sort,
} from '@influxdata/clockface'
import TabbedPageHeader from '../../shared/components/TabbedPageHeader'
import FilterList from '../../shared/components/FilterList'
import NotificationEndpointCards from './NotificationEndpointCards'

// Hooks
import {
  NotificationEndpointsProvider,
  useNotificationEndpoints,
} from './useNotificationEndpoints'

// Types
import {SortKey, SortTypes} from '../../types/sort'
import {NotificationEndpoint} from '../../client'
import NotificationEndpointExplainer from './NotificationEndpointExplainer'

const NotificationEndpointIndex: React.FC = () => {
  const {endpoints} = useNotificationEndpoints()
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Ascending,
  })

  const leftHeader = (
    <>
      <SearchWidget
        search={search}
        placeholder={'Filter Endpoints...'}
        onSearch={setSearch}
      />

      <ResourceSortDropdown
        sortKey={sortOption.key}
        sortType={sortOption.type}
        sortDirection={sortOption.direction}
        onSelect={(sk, sd, st) => {
          setSortOption({
            key: sk,
            type: st,
            direction: sd,
          })
        }}
      />
    </>
  )

  const rightHeader = (
    <Button
      text={'Create Notification Endpoint'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Notification Endpoint'}
      status={ComponentStatus.Default}
    />
  )

  return (
    <>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<NotificationEndpoint>
        list={endpoints}
        search={search}
        searchKeys={['name', 'desc']}
      >
        {filtered => (
          <Grid>
            <Grid.Row>
              <Grid.Column
                widthXS={Columns.Twelve}
                widthSM={filtered.length !== 0 ? Columns.Eight : Columns.Twelve}
                widthMD={filtered.length !== 0 ? Columns.Ten : Columns.Twelve}
              >
                <NotificationEndpointCards
                  search={search}
                  endpoints={filtered}
                  sortKey={sortOption.key}
                  sortType={sortOption.type}
                  sortDirection={sortOption.direction}
                />
              </Grid.Column>

              {filtered.length !== 0 && (
                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthSM={Columns.Four}
                  widthMD={Columns.Two}
                >
                  <NotificationEndpointExplainer />
                </Grid.Column>
              )}
            </Grid.Row>
          </Grid>
        )}
      </FilterList>
    </>
  )
}

export default () => (
  <NotificationEndpointsProvider>
    <NotificationEndpointIndex />
  </NotificationEndpointsProvider>
)
