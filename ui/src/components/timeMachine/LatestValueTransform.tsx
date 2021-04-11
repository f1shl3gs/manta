// Libraries
import React from 'react'

// Components
import {Table} from '@influxdata/giraffe'
import EmptyGraphMessage from 'shared/components/EmptyGraphMessage'

interface Props {
  table: Table
  children: (latestValue: number) => JSX.Element
}

// todo: implement it
const LatestValueTransform: React.FC<Props> = props => {
  const {children, table} = props
  const valueColData = table.getColumn('value', 'number') as number[]
  if (valueColData.length === 0) {
    return <EmptyGraphMessage message={'No latest value found'} />
  }

  return children(valueColData[valueColData.length - 1])
}

export default LatestValueTransform
