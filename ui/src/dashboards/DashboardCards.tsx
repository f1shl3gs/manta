import React, {FunctionComponent} from 'react'
import {useResources} from 'shared/components/GetResources'
import FilterList from 'shared/components/FilterList'
import {Dashboard} from 'types/Dashboard'
import EmptyDashboards from './EmptyDashboards'
import {getSortedResources} from 'shared/utils/sort'
import DashboardCard from './DashboardCard'
import {SortOption} from 'types/Sort'

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
                  onDelete={reload}
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
