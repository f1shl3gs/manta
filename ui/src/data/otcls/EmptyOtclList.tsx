import React from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  searchTerm: string
}

const EmptyOtclList: React.FC<Props> = props => {
  const {searchTerm} = props

  if (searchTerm) {
    return (
      <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
        <EmptyState.Text>
          No <b>Checks</b> match your search
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
      <EmptyState.Text>
        Looks like you have not created a <b>OTCL</b> yet
        <br />
        <br />
        You will need one to collect metrics, traces and logs for your services.
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyOtclList
