// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import ReactGridLayout, {WidthProvider, Layout} from 'react-grid-layout'

// Components
import {Cell} from 'src/types/dashboard'
import GradientBorder from 'src/cells/GradientBorder'
import CellComponent from 'src/cells/Cell'

// Hooks
import {useDashboard} from 'src/dashboards/useDashboard'

const Grid = WidthProvider(ReactGridLayout)

const DASHBOARD_LAYOUT_ROW_HEIGHT = 83.5
const LAYOUT_MARGIN = 4

const layout = (cells: Cell[]): Layout[] => {
  return cells.map(c => {
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

const resizeEventHandler = () => {
  /* void */
}

const Cells: FunctionComponent = () => {
  const {cells, onLayoutChange} = useDashboard()

  useEffect(() => {
    window.addEventListener('resize', resizeEventHandler)

    return () => {
      window.removeEventListener('resize', resizeEventHandler)
    }
  })

  return (
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
