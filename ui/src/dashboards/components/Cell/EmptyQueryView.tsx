// Libraries
import React, {FunctionComponent} from 'react'

import {RemoteDataState} from '@influxdata/clockface'
import {DashboardQuery} from 'src/types/dashboard'
import EmptyGraphMessage from 'src/dashboards/components/Cell/EmptyGraphMessage'
import EmptyGraphError from 'src/dashboards/components/Cell/EmptyGraphError'

interface Props {
  errorMessage?: string
  loading: RemoteDataState
  hasResults: boolean
  isInitialFetch?: boolean
  queries?: DashboardQuery[]
  fallbackNote?: string
  children: JSX.Element | JSX.Element[]
}

const emptyGraphCopy = "Looks like you don't have any queries"

const EmptyQueryView: FunctionComponent<Props> = props => {
  const {loading, queries, errorMessage, hasResults, isInitialFetch} = props

  if (loading === RemoteDataState.NotStarted || (queries && !queries.length)) {
    return (
      <EmptyGraphMessage
        message={emptyGraphCopy}
        testID="empty-graph--no-queries"
      />
    )
  }

  if (errorMessage) {
    return (
      <EmptyGraphError message={errorMessage} testID="empty-graph--error" />
    )
  }

  if ((isInitialFetch || !hasResults) && loading === RemoteDataState.Loading) {
    return <EmptyGraphMessage message="" />
  }

  if (!hasResults) {
    return (
      <EmptyGraphMessage
        message="No results"
        testID={'empty-graph--no-results'}
      />
    )
  }

  return <>{props.children}</>
}

export default EmptyQueryView
