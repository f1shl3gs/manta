// Libraries
import React, {FunctionComponent} from 'react'

// Components
import ErrorBoundary from 'src/shared/components/ErrorBoundary'
import Header from 'src/cells/components/Header'
import Context from 'src/cells/components/Context'
import EmptyGraphMessage from 'src/cells/components/EmptyGraphMessage'
import TimeSeries from 'src/shared/components/TimeSeries'

// Types
import {Cell} from 'src/types/cells'

interface Props {
  cell: Cell
}

const CellComponent: FunctionComponent<Props> = ({cell}) => {
  const {viewProperties} = cell

  const view = (): JSX.Element => {
    if (!viewProperties || viewProperties.queries.length === 0) {
      return (
        <EmptyGraphMessage
          message={'No queries'}
          testID={'empty-graph-message--no-queries'}
        />
      )
    }

    return <TimeSeries viewProperties={viewProperties} />
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
