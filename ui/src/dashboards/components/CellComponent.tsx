// Libraries
import React from 'react';

// Components
import CellHeader from './CellHeader';
import CellContext from './CellContext';
import TimeSeries from './TimeSeries';
import { ViewPropertiesProvider } from 'shared/useViewProperties';

// Types
import { Cell } from 'types/Dashboard';
import ErrorBoundary from '../../shared/components/ErrorBoundary';

interface Props {
  cell: Cell
}

const CellComponent: React.FC<Props> = ({ cell }) => {
  return (
    <>
      <CellHeader name={cell.name || 'Name this Cell'} note={''}>
        <CellContext cell={cell} view={cell.viewProperties} />
      </CellHeader>

      <div className="cell--view">
        {/*<TimeMachine viewProperties={cell.viewProperties} />*/}
        <ViewPropertiesProvider viewProperties={cell.viewProperties}>
          <ErrorBoundary>
            <TimeSeries />
          </ErrorBoundary>
        </ViewPropertiesProvider>
      </div>
    </>
  );
};

export default CellComponent;
