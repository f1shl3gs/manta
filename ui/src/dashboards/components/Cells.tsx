import React, { useEffect, useLayoutEffect } from 'react';
import ReactGridLayout, { WidthProvider, Layout } from 'react-grid-layout';
import { useCells } from '../state/dashboard';
import { Cell } from '../../types';
import CellComponent from './CellComponent';
import GradientBorder from './GradientBorder';

const Grid = WidthProvider(ReactGridLayout);

const DASHBOARD_LAYOUT_ROW_HEIGHT = 83.5;
const LAYOUT_MARGIN = 4;

const layout = (cells: Cell[]): Layout[] => {
  return cells.map((c) => {
    return {
      i: c.id,
      w: c.w || 3,
      h: c.h || 3,
      x: c.x || 0,
      y: c.y || 0
    };
  });
};

const eventHandler = () => {
  console.log('resized event received');
};

const Cells: React.FC = () => {
  const { cells, setCells } = useCells();

  /*useEffect(() => {
    window.addEventListener('resize', eventHandler);
    return () => {
      window.removeEventListener('resize', eventHandler);
    };
  });*/

  useLayoutEffect(() => {
    window.addEventListener('resize', eventHandler);
    return () => {
      window.removeEventListener('resize', eventHandler);
    };
  }, [])

  return (
    <>
      {/* ScrollDetector!? */}
      <Grid
        cols={12}
        layout={layout(cells)}
        rowHeight={DASHBOARD_LAYOUT_ROW_HEIGHT}
        useCSSTransforms={false}
        containerPadding={[0, 0]}
        margin={[LAYOUT_MARGIN, LAYOUT_MARGIN]}
        onLayoutChange={(next) => setCells(next)}
        draggableHandle={'.cell--draggable'}
        isDraggable
        isResizable
      >
        {cells.map((cell) => (
          <div key={cell.id} className="cell">
            <CellComponent cell={cell} />
            <GradientBorder />
          </div>
        ))}
      </Grid>
    </>
  );
};

export default Cells;
