import React from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  searchTerm: string
}

const EmptyVariableList: React.FC<Props> = props => {
  const {searchTerm} = props

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
        Looks like you have not created a <b>Variables</b> yet
        <br />
        <br />
        You will need one to TODO
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyVariableList
