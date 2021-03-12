// Libraries
import React from 'react'

// Components
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  search: string
}

const EmptyChecksList: React.FC<Props> = props => {
  const {search} = props

  if (search) {
    return (
      <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
        <EmptyState.Text>
          No <b>checks</b> match your search
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
      <EmptyState.Text>
        Looks like you have not created a <b>Check</b> yet
        <br />
        <br />
        You will need one to be notified about
        <br />
        any changes in system status
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyChecksList
