import React from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  search: string
}

const EmptyEndpointsList: React.FC<Props> = props => {
  const {search} = props

  if (search) {
    return (
      <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
        <EmptyState.Text>
          No <b>notification endpoints</b> match your search
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
      <EmptyState.Text>
        Looks like you have not created a <b>Notification Endpoint</b> yet
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyEndpointsList
