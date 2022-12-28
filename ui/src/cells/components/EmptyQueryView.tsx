// Libraries
import React, {FunctionComponent} from 'react'

import {RemoteDataState} from '@influxdata/clockface'
import {DashboardQuery} from 'src/types/dashboards'
import EmptyGraphMessage from 'src/cells/components/EmptyGraphMessage'
import EmptyGraphError from 'src/cells/components/EmptyGraphError'

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

const EmptyQueryView: FunctionComponent<Props> = ({
  loading,
  queries,
  errorMessage,
  hasResults,
  isInitialFetch,
  children
}) => {
  if (!hasResults) {
    return (
      <EmptyGraphMessage
        message="No results"
        testID={'empty-graph--no-results'}
      />
      )
  }

  if (loading === RemoteDataState.NotStarted || (queries && queries.length === 0)) {
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

  return <>{children}</>
}

export default EmptyQueryView
