// Libraries
import React, {useState} from 'react'

// Hooks
import {useVariables, VariablesProvider} from './useVariables'

// Utils
import withProvider from '../../utils/withProvider'
import SearchWidget from '../../shared/components/SearchWidget'
import {SortKey, SortTypes} from '../../types/sort'
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
import ResourceSortDropdown from '../../shared/components/ResourceSortDropdown'
import {useHistory} from 'react-router-dom'
import {useOrgID} from '../../shared/useOrg'
import TabbedPageHeader from '../../shared/components/TabbedPageHeader'
import {Variable} from '../../types/Variable'
import FilterList from '../../shared/components/FilterList'
import VariableCards from './VariableCards'

const Variables: React.FC = () => {
  const history = useHistory()
  const orgID = useOrgID()
  const [searchTerm, setSearchTerm] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Descending,
  })
  const {variables, loading} = useVariables()

  const leftHeader = (
    <>
      <SearchWidget
        search={searchTerm}
        placeholder={'Filter Variables'}
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
      text={'Create Variable'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Variable'}
      status={ComponentStatus.Default}
      onClick={() => {
        history.push(`/orgs/${orgID}/settings/variables/new`)
      }}
    />
  )

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<Variable>
        list={variables}
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
                <VariableCards
                  list={filtered}
                  searchTerm={searchTerm}
                  sortKey={sortOption.key}
                  sortType={sortOption.type}
                  sortDirection={sortOption.direction}
                />
              </Grid.Column>
            </Grid.Row>
          </Grid>
        )}
      </FilterList>
    </SpinnerContainer>
  )
}

export default withProvider(VariablesProvider, Variables)
