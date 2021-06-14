import React from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  searchTerm: string
  resourceName: string
  description?: string
}

const EmptyResourceList: React.FC<Props> = props => {
  const {searchTerm, resourceName, description} = props

  if (searchTerm) {
    return (
      <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
        <EmptyState.Text>
          No <b>Variables</b> match your search
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
      <EmptyState.Text>
        Looks like you have not created a <b>{resourceName}</b> yet
        <br />
        <br />
        {description}
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyResourceList
