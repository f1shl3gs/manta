// Libraries
import React, {FunctionComponent, lazy, useState} from 'react'
import {Route, Routes} from 'react-router-dom'

// Components
import {Columns, Grid, Sort} from '@influxdata/clockface'
import FilterList from 'src/shared/components/FilterList'
import {AutoSizer} from 'react-virtualized'
import ResourceSortDropdown from 'src/shared/components/resource_sort_dropdown/ResourceSortDropdown'
import SearchWidget from 'src/shared/components/SearchWidget'
import CreateConfigurationButton from 'src/configurations/components/CreateConfigurationButton'
import ConfigurationCard from 'src/configurations/components/ConfigurationCard'
import {getSortedResources} from 'src/shared/utils/sort'
import ConfigurationExplainer from 'src/configurations/components/ConfigurationExplainer'

// Types
import {SortTypes} from 'src/types/sort'
import {Configuration} from 'src/types/configuration'
import {useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'
import {getAll} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import GetResources from 'src/resources/components/GetResources'
import EmptyResources from '../../resources/components/EmptyResources'

const ConfigurationWizard = lazy(
  () => import('src/configurations/components/ConfigurationWizard')
)

const DEFAULT_PAGINATION_CONTROL_HEIGHT = 62
const DEFAULT_TAB_NAVIGATION_HEIGHT = 62

const ConfigurationIndex: FunctionComponent = () => {
  const configs = useSelector((state: AppState) => {
    return getAll<Configuration>(state, ResourceType.Configurations)
  })
  const [search, setSearch] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated',
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
        resource={ResourceType.Configurations}
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
                      list={configs}
                      search={search}
                      searchKeys={['name', 'desc']}
                    >
                      {filtered => {
                        if (filtered && filtered.length === 0) {
                          return (
                            <EmptyResources
                              searchTerm={search}
                              resource={ResourceType.Configurations}
                              createButton={<CreateConfigurationButton />}
                            />
                          )
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

export default () => (
  <GetResources resources={[ResourceType.Configurations]}>
    <ConfigurationIndex />
  </GetResources>
)
