import React, {FunctionComponent} from 'react'
import {useResources} from 'src/shared/components/GetResources'
import FilterList from 'src/shared/components/FilterList'
import {Dashboard} from 'src/types/dashboard'
import EmptyDashboards from 'src/dashboards/EmptyDashboards'
import {getSortedResources} from 'src/utils/sort'
import DashboardCard from 'src/dashboards/DashboardCard'
import {SortOption} from 'src/types/sort'

interface Props {
  search: string
  sortOption: SortOption
}

const DashboardCards: FunctionComponent<Props> = props => {
  const {sortOption, search} = props
  const {resources, reload} = useResources()

  return (
    <FilterList<Dashboard>
      list={resources}
      search={search}
      searchKeys={['name', 'desc']}
    >
      {filtered => {
        if (filtered && filtered.length === 0) {
          return <EmptyDashboards searchTerm={search} />
        }

        return (
          <div style={{height: '100%', display: 'grid'}}>
            <div className={'dashboards-card-grid'}>
              {getSortedResources<Dashboard>(
                filtered,
                sortOption.key,
                sortOption.type,
                sortOption.direction
              ).map(dashboard => (
                <DashboardCard
                  key={dashboard.id}
                  dashboard={dashboard}
                  reload={reload}
                />
              ))}
            </div>
          </div>
        )
      }}
    </FilterList>
  )
}

export default DashboardCards
