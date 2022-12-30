// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Component
import FilterList from 'src/shared/components/FilterList'
import DashboardCard from 'src/dashboards/components/DashboardCard'
import EmptyResources from 'src/resources/components/EmptyResources'
import CreateDashboardButton from 'src/dashboards/components/CreateDashboardButton'

// Types
import {ResourceType} from 'src/types/resources'
import {Dashboard} from 'src/types/dashboards'
import {getAll} from 'src/resources/selectors'
import {AppState} from 'src/types/stores'

// Utils
import {getSortedResources} from 'src/shared/utils/sort'

const mstp = (state: AppState) => {
  const dashboards = getAll<Dashboard>(state, ResourceType.Dashboards)
  const {sortOptions, searchTerm} = state.resources[ResourceType.Dashboards]

  return {
    dashboards,
    searchTerm,
    sortOptions,
  }
}

const connector = connect(mstp, null)
type Props = ConnectedProps<typeof connector>

const DashboardCards: FunctionComponent<Props> = ({
  dashboards,
  sortOptions,
  searchTerm,
}) => {
  return (
    <FilterList<Dashboard>
      list={dashboards}
      search={searchTerm}
      searchKeys={['name', 'desc']}
    >
      {filtered => {
        if (filtered && filtered.length === 0) {
          return (
            <EmptyResources
              resource={ResourceType.Dashboards}
              searchTerm={searchTerm}
              createButton={<CreateDashboardButton />}
            />
          )
        }

        return (
          <div style={{height: '100%', display: 'grid'}}>
            <div className={'dashboards-card-grid'}>
              {getSortedResources<Dashboard>(
                filtered,
                sortOptions.key,
                sortOptions.type,
                sortOptions.direction
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

export default connector(DashboardCards)
