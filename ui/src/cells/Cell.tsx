import React, {FunctionComponent} from 'react'

import {Cell} from 'src/types/dashboard'
import ErrorBoundary from 'src/shared/components/ErrorBoundary'
import Header from 'src/cells/Header'
import Context from 'src/cells/Context'
import EmptyGraphMessage from 'src/cells/EmptyGraphMessage'
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
