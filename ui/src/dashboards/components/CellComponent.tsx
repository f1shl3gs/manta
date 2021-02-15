// Libraries
import React from 'react'

// Components
import CellHeader from './CellHeader'
import CellContext from './CellContext'
import TimeSeries from './TimeSeries'
import ErrorBoundary from 'shared/components/ErrorBoundary'

// Hooks
import {ViewPropertiesProvider} from 'shared/useViewProperties'

// Types
import {Cell} from 'types/Dashboard'

interface Props {
  cell: Cell
}

const CellComponent: React.FC<Props> = ({cell}) => {
  return (
    <ErrorBoundary>
      <CellHeader name={cell.name || 'Name this Cell'} note={''}>
        <CellContext cell={cell} view={cell.viewProperties} />
      </CellHeader>

      <div className="cell--view">
        <ViewPropertiesProvider viewProperties={cell.viewProperties}>
          <TimeSeries />
        </ViewPropertiesProvider>
      </div>
    </ErrorBoundary>
  )
}

export default CellComponent
