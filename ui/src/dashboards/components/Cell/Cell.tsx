import React, {FunctionComponent} from 'react'

import {Cell} from 'src/types/Dashboard'
import ErrorBoundary from 'src/shared/components/ErrorBoundary'
import Header from 'src/dashboards/components/Cell/Header'
import Context from 'src/dashboards/components/Cell/Context'
import EmptyGraphMessage from 'src/dashboards/components/Cell/EmptyGraphMessage'
import {ViewPropertiesProvider} from 'src/visualization/TimeMachine/useViewProperties'
import TimeSeries from 'src/shared/components/TimeSeries'

interface Props {
  cell: Cell
}

const CellComponent: FunctionComponent<Props> = ({cell}) => {
  const {viewProperties} = cell

  const view = (): JSX.Element => {
    if (!viewProperties) {
      return (
        <EmptyGraphMessage
          message={'No queries'}
          testID={'empty-graph-message--no-queries'}
        />
      )
    }

    return (
      <ViewPropertiesProvider viewProperties={viewProperties}>
        <TimeSeries viewProperties={viewProperties} />
      </ViewPropertiesProvider>
    )
  }

  return (
    <ErrorBoundary>
      <Header name={cell.name || 'Name this Cell'} note={''}>
        <Context cell={cell} />
      </Header>

      <div className={'cell--view'}>{view()}</div>
    </ErrorBoundary>
  )
}

export default CellComponent
