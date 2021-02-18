// Libraries
import React from 'react'

// Components
import {ResourceList} from '@influxdata/clockface'
import EmptyChecksList from './EmptyChecksList'
import CheckCard from './CheckCard'

// Types
import {Check} from '../types/Check'

interface Props {
  search: string
  checks: Check[]
}

const CheckCards: React.FC<Props> = props => {
  const {search, checks} = props

  const body = (filtered: Check[]) => (
    <ResourceList.Body emptyState={<EmptyChecksList search={search} />}>
      {filtered.map(check => (
        <CheckCard check={check} />
      ))}
    </ResourceList.Body>
  )

  return <ResourceList>{body(checks)}</ResourceList>
}

export default CheckCards
