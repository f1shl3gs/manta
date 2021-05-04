// Libraries
import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  Button,
  Columns,
  ComponentColor,
  ComponentStatus,
  Grid,
  IconFont,
  Sort,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import TabbedPageHeader from '../../shared/components/TabbedPageHeader'
import SearchWidget from '../../shared/components/SearchWidget'
import ResourceSortDropdown from '../../shared/components/ResourceSortDropdown'
import FilterList from '../../shared/components/FilterList'
import ScraperCards from './ScraperCards'
import ScraperExplainer from './ScraperExplainer'

// Hooks
import {useOrgID} from '../../shared/useOrg'
import {Scraper} from '../../types/scrapers'

// Utils
import withProvider from '../../utils/withProvider'
import {ScrapersProvider, useScrapers} from './useScrapers'

// Types
import {SortKey, SortTypes} from '../../types/sort'

const Scrapers: React.FC = () => {
  const {scrapers, loading} = useScrapers()
  const history = useHistory()
  const orgID = useOrgID()
  const [searchTerm, setSearchTerm] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  const leftHeader = (
    <>
      <SearchWidget
        search={searchTerm}
        placeholder={'Filter Scrapers...'}
        onSearch={setSearchTerm}
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
      text={'Create Scraper'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Scraper'}
      status={ComponentStatus.Default}
      onClick={() => {
        history.push(`/orgs/${orgID}/data/scrapers/new`)
      }}
    />
  )

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<Scraper>
        list={scrapers}
        search={searchTerm}
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
                <ScraperCards
                  list={filtered}
                  searchTerm={searchTerm}
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
                  <ScraperExplainer />
                </Grid.Column>
              )}
            </Grid.Row>
          </Grid>
        )}
      </FilterList>
    </SpinnerContainer>
  )
}

export default withProvider(ScrapersProvider, Scrapers)
