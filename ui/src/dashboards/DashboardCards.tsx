// Libraries
import React, {useCallback} from 'react'
import moment from 'moment'

// Components
import DashboardCard from './components/DashboardCard'

// Hooks
import {useDashboards} from './useDashboards'
import {useFetch} from 'shared/useFetch'
import {SortKey, SortTypes} from '../types/sort'
import {Sort} from '@influxdata/clockface'

import {Dashboard} from '../types/Dashboard'

// Utils
import {getSortedResources} from 'utils/sort'

interface Props {
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

const DashboardCards: React.FC<Props> = props => {
  const {sortKey, sortType, sortDirection} = props
  const {dashboards, refresh} = useDashboards()

  const {del} = useFetch(`/api/v1/dashboards`, {})
  const onDeleteDashboard = useCallback(
    (id: string) => {
      del(id)
        .then(() => {
          refresh()
        })
        .catch(err => {
          console.log('delete dashboard err', err)
        })
    },
    [del, refresh]
  )

  const body = (filtered: Dashboard[]) =>
    getSortedResources<Dashboard>(
      filtered,
      sortKey,
      sortType,
      sortDirection
    ).map(d => (
      <DashboardCard
        key={d.id}
        id={d.id}
        name={d.name}
        desc={d.desc}
        updatedAt={moment(d.updated).fromNow()}
        onDeleteDashboard={onDeleteDashboard}
      />
    ))

  return (
    <div style={{height: '100%', display: 'grid'}}>
      <div className={'dashboards-card-grid'}>{body(dashboards)}</div>
    </div>
  )
}

export default DashboardCards
