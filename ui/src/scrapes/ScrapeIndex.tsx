// Libraries
import React, {FunctionComponent, useState} from 'react'

// Components
import {Columns, Grid, Sort} from '@influxdata/clockface'
import FilterList from 'src/shared/components/FilterList'
import {getSortedResources} from 'src/utils/sort'
import ScrapeCard from 'src/scrapes/ScrapeCard'
import EmptyScrapes from 'src/scrapes/EmptyScrapes'
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import CreateScrapeButton from 'src/scrapes/CreateScrapeButton'
import ScrapeExplainer from 'src/scrapes/ScrapeExplainer'

// Types
import {ResourceType} from 'src/types/resources'
import {Scrape} from 'src/types/scrape'
import {SortKey, SortTypes} from 'src/types/sort'

// Hooks
import {useSelector} from 'react-redux'

// Actions
import GetResources from 'src/resources/components/GetResources'
import {AppState} from 'src/types/stores'
import {getAll} from 'src/resources/selectors'

const ScrapeIndex: FunctionComponent = () => {
  const scrapes = useSelector((state: AppState) =>
    getAll<Scrape>(state, ResourceType.Scrapes)
  )
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
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
      <div className={'tabbed-page--header'}>
        {left}
        {right}
      </div>

      <Grid>
        <Grid.Row>
          <Grid.Column
            widthXS={Columns.Twelve}
            widthSM={Columns.Eight}
            widthMD={Columns.Nine}
            widthLG={Columns.Ten}
          >
            <FilterList<Scrape>
              list={scrapes}
              search={search}
              searchKeys={['name', 'desc']}
            >
              {filtered => {
                if (filtered && filtered.length === 0) {
                  return <EmptyScrapes searchTerm={search} />
                }

                return (
                  <div style={{height: '100%', display: 'grid'}}>
                    {getSortedResources<Scrape>(
                      filtered,
                      sortOption.key,
                      sortOption.type,
                      sortOption.direction
                    ).map(sc => (
                      <ScrapeCard key={sc.id} scrape={sc} />
                    ))}
                  </div>
                )
              }}
            </FilterList>
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

export default () => (
  <GetResources resources={[ResourceType.Scrapes]}>
    <ScrapeIndex />
  </GetResources>
)
