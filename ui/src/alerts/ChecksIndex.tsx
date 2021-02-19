// Libraries
import React, {useState} from 'react'

// Components
import {
  Button,
  Columns,
  ComponentColor,
  ComponentStatus,
  Grid,
  IconFont,
  Sort,
} from '@influxdata/clockface'
import TabbedPageHeader from 'shared/components/TabbedPageHeader'
import SearchWidget from 'shared/components/SearchWidget'
import FilterList from 'shared/components/FilterList'
import ResourceSortDropdown from '../shared/components/ResourceSortDropdown'

// Hooks
import {useChecks} from './useChecks'

// Types
import {Check} from 'types/Check'
import {SortKey, SortTypes} from 'types/sort'
import CheckCards from './CheckCards'
import CheckExplainer from './CheckExplainer'

const ChecksIndex: React.FC = () => {
  const {checks, remoteDataState} = useChecks()
  const title = `Checks`
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Ascending,
  })

  const tooltipContents = (
    <>
      A <strong>Check</strong> is a periodic query that the system
      <br />
      performs against your time series data
      <br />
      that will generate a status
      <br />
      <br />
    </>
  )

  const leftHeader = (
    <>
      <SearchWidget
        search={search}
        placeholder={'Filter Checks...'}
        onSearch={v => setSearch(v)}
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
      text={'Create Check'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Check'}
      status={ComponentStatus.Default}
    />
  )

  return (
    <>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<Check>
        list={checks}
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
                <CheckCards
                  search={search}
                  checks={filtered}
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
                  <CheckExplainer />
                </Grid.Column>
              )}
            </Grid.Row>
          </Grid>
        )}
      </FilterList>
    </>
  )
}

export default ChecksIndex
