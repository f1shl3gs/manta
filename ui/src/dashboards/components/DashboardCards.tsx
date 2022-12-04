// Libraries
import React, {FunctionComponent} from 'react'

// Component
import FilterList from 'src/shared/components/FilterList'
import EmptyDashboards from 'src/dashboards/components/EmptyDashboards'
import DashboardCard from 'src/dashboards/components/DashboardCard'

// Hooks
import {useSelector} from 'react-redux'

// Types
import {ResourceType} from 'src/types/resources'
import {Dashboard, DashboardSortParams} from 'src/types/dashboards'
import {getAll} from 'src/resources/selectors'
import {AppState} from 'src/types/stores'

// Utils
import {getSortedResources} from 'src/utils/sort'

interface Props {
  search: string
  sortOption: DashboardSortParams
}

const DashboardCards: FunctionComponent<Props> = ({sortOption, search}) => {
  const dashboards = useSelector((state: AppState) =>
    getAll<Dashboard>(state, ResourceType.Dashboards)
  )

  return (
    <FilterList<Dashboard>
      list={dashboards}
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
                <DashboardCard key={dashboard.id} id={dashboard.id} />
              ))}
            </div>
          </div>
        )
      }}
    </FilterList>
  )
}

export default DashboardCards
