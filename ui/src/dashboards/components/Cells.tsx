import React, {useEffect} from 'react'
import ReactGridLayout, {WidthProvider, Layout} from 'react-grid-layout'
import {Cell} from 'types/Dashboard'

import CellComponent from './CellComponent'
import GradientBorder from './GradientBorder'
import {useDashboard} from './useDashboard'

const Grid = WidthProvider(ReactGridLayout)

const DASHBOARD_LAYOUT_ROW_HEIGHT = 83.5
const LAYOUT_MARGIN = 4

const layout = (cells: Cell[]): Layout[] => {
  return cells.map((c) => {
    return {
      i: c.id,
      w: c.w || 0,
      h: c.h || 0,
      x: c.x || 0,
      y: c.y || 0,
      static: false,
      moved: false,
      isBounded: undefined,
      isDraggable: undefined,
      isResizable: undefined,
      maxH: undefined,
      maxW: undefined,
      minH: undefined,
      minW: undefined,
      resizeHandles: undefined,
    }
  })
}

const eventHandler = () => {
  console.log('resized event received')
}

const Cells: React.FC = () => {
  const {cells, onLayoutChange} = useDashboard()

  // useEffect(() => {
  //   window.addEventListener('resize', eventHandler);
  //   return () => {
  //     window.removeEventListener('resize', eventHandler);
  //   };
  // });

  return (
    /* <ScrollDetector component="dashboard" /> */
    <Grid
      cols={12}
      layout={layout(cells)}
      rowHeight={DASHBOARD_LAYOUT_ROW_HEIGHT}
      useCSSTransforms={false}
      containerPadding={[0, 0]}
      margin={[LAYOUT_MARGIN, LAYOUT_MARGIN]}
      onLayoutChange={onLayoutChange}
      draggableHandle={'.cell--draggable'}
      isDraggable
      isResizable
      measureBeforeMount={true}
    >
      {cells.map((cell: Cell) => (
        <div key={cell.id} className="cell">
          <CellComponent cell={cell} />
          <GradientBorder />
        </div>
      ))}
    </Grid>
  )
}

export default Cells
