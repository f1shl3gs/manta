// Libraries
import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

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
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import TabbedPageHeader from '../../shared/components/TabbedPageHeader'
import FilterList from '../../shared/components/FilterList'
import {useOtcls, OtclsProvider} from './useOtcls'
import OtclCards from './OtclCards'
import OtclExplainer from './OtclExplainer'

// Hooks
import {useOrgID} from '../../shared/useOrg'

// Types
import {SortKey, SortTypes} from '../../types/sort'
import {Otcl} from '../../types/otcl'
import withProvider from '../../utils/withProvider'

const Otcls: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Descending,
  })
  const {otcls, loading, onDelete} = useOtcls()
  const history = useHistory()
  const orgID = useOrgID()

  const leftHeader = (
    <>
      <SearchWidget
        search={searchTerm}
        placeholder={'Filter Otcls...'}
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
      text={'Create Otcl'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Otcl configuration'}
      status={ComponentStatus.Default}
      onClick={() => {
        history.push(`/orgs/${orgID}/data/otcls/new`)
      }}
    />
  )

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<Otcl>
        list={otcls}
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
                <OtclCards
                  list={filtered}
                  searchTerm={searchTerm}
                  sortKey={sortOption.key}
                  sortType={sortOption.type}
                  sortDirection={sortOption.direction}
                  onDelete={onDelete}
                />
              </Grid.Column>

              {filtered.length !== 0 && (
                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthSM={Columns.Four}
                  widthMD={Columns.Two}
                >
                  <OtclExplainer />
                </Grid.Column>
              )}
            </Grid.Row>
          </Grid>
        )}
      </FilterList>
    </SpinnerContainer>
  )
}

export default withProvider(OtclsProvider, Otcls)
