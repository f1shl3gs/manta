// Libraries
import React, {FunctionComponent, useState} from 'react'

// Components
import {Columns, Grid, Sort} from '@influxdata/clockface'
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'
import CreateScrapeButton from 'src/scrapes/components/CreateScrapeButton'
import ScrapeExplainer from 'src/scrapes/components/ScrapeExplainer'

// Types
import {ResourceType} from 'src/types/resources'
import {SortTypes} from 'src/types/sort'

// Actions
import GetResources from 'src/resources/components/GetResources'
import TabbedPageHeader from 'src/shared/components/TabbedPageHeader'
import ScrapeList from './ScrapeList'

const ScrapeIndex: FunctionComponent = () => {
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated',
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  const left = (
    <div className={'tabbed-page--header-left'}>
      <SearchWidget
        search={search}
        placeholder={'Filter Scrapes'}
        onSearch={s => setSearch(s)}
      />

      <ResourceSortDropdown
        resource={ResourceType.Scrapes}
        sortKey={sortOption.key}
        sortType={sortOption.type}
        sortDirection={sortOption.direction}
        onSelect={(key, direction, type) => {
          setSortOption({key, type, direction})
        }}
      />
    </div>
  )

  const right = (
    <div className={'tabbed-page--header-right'}>
      <CreateScrapeButton />
    </div>
  )

  return (
    <>
      <TabbedPageHeader left={left} right={right} />

      <Grid>
        <Grid.Row>
          <Grid.Column
            widthXS={Columns.Twelve}
            widthSM={Columns.Eight}
            widthMD={Columns.Nine}
            widthLG={Columns.Ten}
          >
            <GetResources resources={[ResourceType.Scrapes]}>
              <ScrapeList search={search} sortOption={sortOption} />
            </GetResources>
          </Grid.Column>

          <Grid.Column
            widthXS={Columns.Twelve}
            widthSM={Columns.Four}
            widthMD={Columns.Three}
            widthLG={Columns.Two}
          >
            <ScrapeExplainer />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </>
  )
}

export default ScrapeIndex
