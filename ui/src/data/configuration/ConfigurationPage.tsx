// Libraries
import React, {FunctionComponent, lazy, useState} from 'react'
import {Route, Routes} from 'react-router-dom'

// Components
import {Columns, Grid, Sort} from '@influxdata/clockface'
import FilterList from 'src/shared/components/FilterList'
import {AutoSizer} from 'react-virtualized'
import {
  ResourceType,
  useResources,
  withResources,
} from 'src/shared/components/GetResources'
import ResourceSortDropdown from 'src/shared/components/ResourceSortDropdown'
import SearchWidget from 'src/shared/components/SearchWidget'
import CreateConfigurationButton from 'src/data/configuration/CreateConfigurationButton'
import EmptyConfigurations from 'src/data/configuration/EmptyConfigurations'
import ConfigurationCard from 'src/data/configuration/ConfigurationCard'
import {getSortedResources} from 'src/utils/sort'
import ConfigurationExplainer from 'src/data/configuration/ConfigurationExplainer'

// Types
import {SortKey, SortTypes} from 'src/types/sort'
import {Configuration} from 'src/types/configuration'

const ConfigurationWizard = lazy(
  () => import('src/data/configuration/ConfigurationWizard')
)

const DEFAULT_PAGINATION_CONTROL_HEIGHT = 62
const DEFAULT_TAB_NAVIGATION_HEIGHT = 62

const ConfigurationPage: FunctionComponent = () => {
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
        onSearch={t => setSearch(t)}
        placeholder={'Filter Configurations...'}
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
    </div>
  )

  const right = (
    <div className={'tabbed-page--header-right'}>
      <CreateConfigurationButton />
    </div>
  )

  return (
    <>
      <Routes>
        <Route path="new" element={<ConfigurationWizard />} />
      </Routes>

      <AutoSizer>
        {({width, height}) => {
          const heightWithPagination =
            DEFAULT_PAGINATION_CONTROL_HEIGHT + DEFAULT_TAB_NAVIGATION_HEIGHT
          const adjustedHeight = height - heightWithPagination - 60

          return (
            <>
              <div className={'tabbed-page--header'} style={{width}}>
                {left}
                {right}
              </div>

              <Grid style={{width, height: adjustedHeight}}>
                <Grid.Row>
                  <Grid.Column
                    widthXS={Columns.Twelve}
                    widthSM={Columns.Eight}
                    widthMD={Columns.Nine}
                    widthLG={Columns.Ten}
                  >
                    <FilterList<Configuration>
                      list={resources}
                      search={search}
                      searchKeys={['name', 'desc']}
                    >
                      {filtered => {
                        if (filtered && filtered.length === 0) {
                          return <EmptyConfigurations searchTerm={search} />
                        }

                        return (
                          <div style={{height: '100%', display: 'grid'}}>
                            {getSortedResources<Configuration>(
                              filtered,
                              sortOption.key,
                              sortOption.type,
                              sortOption.direction
                            ).map(cf => (
                              <ConfigurationCard
                                key={cf.id}
                                configuration={cf}
                              />
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
                    <ConfigurationExplainer />
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

export default withResources(ConfigurationPage, ResourceType.Configurations)
