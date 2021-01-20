// Libraries
import React, {useCallback} from 'react'
import moment from 'moment'

// Components
import DashboardCard from './components/DashboardCard'

// Hooks
import {useDashboards} from './useDashboards'
import {useFetch} from 'use-http'

const DashboardCards: React.FC = () => {
  const {dashboards, refresh} = useDashboards()

  const {del} = useFetch(`/api/v1/dashboards`, {})
  const onDeleteDashboard = useCallback(
    (id: string) => {
      del(id)
        .then(() => {
          refresh()
        })
        .catch((err) => {
          console.log('delete dashboard err', err)
        })
    },
    [del]
  )

  return (
    <div style={{height: '100%', display: 'grid'}}>
      <div className={'dashboards-card-grid'}>
        {dashboards?.map((d) => (
          <DashboardCard
            key={d.id}
            id={d.id}
            name={d.name}
            desc={d.desc}
            updatedAt={moment(d.updated).fromNow()}
            onDeleteDashboard={onDeleteDashboard}
          />
        ))}
      </div>
    </div>
  )
}

export default DashboardCards
