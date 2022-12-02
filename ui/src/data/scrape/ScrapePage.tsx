import {Columns, Grid, Sort} from '@influxdata/clockface'
import React, {FunctionComponent, useState} from 'react'
import FilterList from 'src/shared/components/FilterList'
import {
  ResourceType,
  useResources,
  withResources,
} from 'src/shared/components/GetResources'
import {getSortedResources} from 'src/utils/sort'
import {Scrape} from 'src/types/scrape'
import {SortKey, SortTypes} from 'src/types/sort'
import ScrapeCard from 'src/data/scrape/ScrapeCard'
import EmptyScrapes from 'src/data/scrape/EmptyScrapes'
import SearchWidget from 'src/shared/components/SearchWidget'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import CreateScrapeButton from 'src/data/scrape/CreateScrapeButton'
import ScrapeExplainer from 'src/data/scrape/ScrapeExplainer'

const ScrapePage: FunctionComponent = () => {
  const {resources} = useResources()
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
              list={resources}
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

export default withResources(ScrapePage, ResourceType.Scrapes)
