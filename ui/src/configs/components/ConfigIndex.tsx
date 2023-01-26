// Libraries
import React, {FunctionComponent, lazy, useState} from 'react'
import {Route, Routes} from 'react-router-dom'

// Components
import {Columns, Grid, Sort} from '@influxdata/clockface'
import FilterList from 'src/shared/components/FilterList'
import {AutoSizer} from 'react-virtualized'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'
import SearchWidget from 'src/shared/components/SearchWidget'
import CreateConfigButton from 'src/configs/components/CreateConfigButton'
import ConfigCard from 'src/configs/components/ConfigCard'
import {getSortedResources} from 'src/shared/utils/sort'
import ConfigExplainer from 'src/configs/components/ConfigExplainer'

// Types
import {SortTypes} from 'src/types/sort'
import {Config} from 'src/types/config'
import {useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'
import {getAll} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import GetResources from 'src/resources/components/GetResources'
import EmptyResources from 'src/resources/components/EmptyResources'
import TabbedPageHeader from 'src/shared/components/TabbedPageHeader'

const ConfigWizard = lazy(() => import('src/configs/components/ConfigWizard'))

const DEFAULT_PAGINATION_CONTROL_HEIGHT = 62
const DEFAULT_TAB_NAVIGATION_HEIGHT = 62

const ConfigIndex: FunctionComponent = () => {
  const configs = useSelector((state: AppState) => {
    return getAll<Config>(state, ResourceType.Configs)
  })
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated',
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  const left = (
    <>
      <SearchWidget
        onSearch={t => setSearch(t)}
        placeholder={'Filter Configs...'}
        search={search}
      />

      <ResourceSortDropdown
        resource={ResourceType.Configs}
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
    </>
  )

  const right = <CreateConfigButton />

  return (
    <>
      <Routes>
        <Route path="new" element={<ConfigWizard />} />
      </Routes>

      <AutoSizer>
        {({width, height}) => {
          const heightWithPagination =
            DEFAULT_PAGINATION_CONTROL_HEIGHT + DEFAULT_TAB_NAVIGATION_HEIGHT
          const adjustedHeight = height - heightWithPagination - 60

          return (
            <>
              <TabbedPageHeader left={left} right={right} style={{width}} />

              <Grid style={{width, height: adjustedHeight}}>
                <Grid.Row>
                  <Grid.Column
                    widthXS={Columns.Twelve}
                    widthSM={Columns.Eight}
                    widthMD={Columns.Nine}
                    widthLG={Columns.Ten}
                  >
                    <FilterList<Config>
                      list={configs}
                      search={search}
                      searchKeys={['name', 'desc']}
                    >
                      {filtered => {
                        if (filtered && filtered.length === 0) {
                          return (
                            <EmptyResources
                              searchTerm={search}
                              resource={ResourceType.Configs}
                              createButton={<CreateConfigButton />}
                            />
                          )
                        }

                        return (
                          <div style={{height: '100%', display: 'grid'}}>
                            {getSortedResources<Config>(
                              filtered,
                              sortOption.key,
                              sortOption.type,
                              sortOption.direction
                            ).map(cf => (
                              <ConfigCard key={cf.id} config={cf} />
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
                    <ConfigExplainer />
                  </Grid.Column>
                </Grid.Row>
              </Grid>
            </>
          )
        }}
      </AutoSizer>
    </>
  )
}

export default () => (
  <GetResources resources={[ResourceType.Configs]}>
    <ConfigIndex />
  </GetResources>
)
