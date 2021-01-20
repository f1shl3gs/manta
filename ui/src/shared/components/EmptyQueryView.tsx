// Libraries
import React from 'react'

// Types
import {DashboardQuery} from 'types/Dashboard'
import {RemoteDataState} from '@influxdata/clockface'
import EmptyGraphMessage from './EmptyGraphMessage'
import EmptyGraphError from './EmptyGraphError'

interface Props {
  loading: RemoteDataState
  errorMessage?: string
  fallbackNote?: string
  queries?: DashboardQuery[]
  hasResults: boolean
}

const EmptyQueryView: React.FC<Props> = (props) => {
  const {loading, queries, errorMessage, fallbackNote, hasResults} = props

  if (queries && !queries.length) {
    return (
      <EmptyGraphMessage
        message={
          'Looks like you don’t have any queries. Be a lot cooler if you did!'
        }
      />
    )
  }

  if (errorMessage !== undefined) {
    return <EmptyGraphError message={errorMessage} />
  }

  if (loading === RemoteDataState.Loading) {
    return <EmptyGraphMessage />
  }

  if (fallbackNote) {
    return <div>{fallbackNote}</div>
  }

  if (!hasResults) {
    return (
      <EmptyGraphMessage
        message="No Results"
        testID="empty-graph--no-results"
      />
    )
  }

  return <>{props.children}</>
}

export default EmptyQueryView
