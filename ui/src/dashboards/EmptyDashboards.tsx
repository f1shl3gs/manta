// Libraries
import React from 'react'

// Components
import {ComponentSize, EmptyState} from '@influxdata/clockface'
import CreateDashboardButton from 'src/dashboards/CreateDashboardButton'

interface Props {
  searchTerm: string
}

const EmptyDashboards: React.FC<Props> = props => {
  const {searchTerm} = props

  if (searchTerm) {
    return (
      <EmptyState
        size={ComponentSize.Large}
        testID={'no-match-dashboards-list'}
      >
        <EmptyState.Text>
          No Dashboards match your search term <b>{searchTerm}</b>
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Large} testID={'empty-dashboards-list'}>
      <EmptyState.Text>
        Looks like you don't have any <b>Dashboards</b>, why not create one?
      </EmptyState.Text>

      <CreateDashboardButton />
    </EmptyState>
  )
}

export default EmptyDashboards
