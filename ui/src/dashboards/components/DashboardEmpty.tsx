// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {ComponentSize, EmptyState} from '@influxdata/clockface'
import CreateCellButton from './CreateCellButton'

const DashboardEmpty: FunctionComponent = () => {
  return (
    <div className={'dashboard-empty'}>
      <EmptyState size={ComponentSize.Large}>
        <EmptyState.Text>
          The Dashboard doesn't have any <b>Cells</b>, let's create some!
        </EmptyState.Text>

        <CreateCellButton />
      </EmptyState>
    </div>
  )
}

export default DashboardEmpty
