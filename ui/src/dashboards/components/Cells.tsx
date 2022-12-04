// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import ReactGridLayout, {Layout, WidthProvider} from 'react-grid-layout'

// Components
import {Cell} from 'src/types/cells'
import GradientBorder from 'src/cells/components/GradientBorder'
import CellComponent from 'src/cells/components/Cell'
import {connect, ConnectedProps, useDispatch} from 'react-redux'
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'
import {updateLayout} from 'src/dashboards/actions/thunks'
import {getCells} from 'src/cells/selectors'

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

const mstp = (state: AppState) => {
  const dashboardID = state.resources[ResourceType.Dashboards].current

  return {
    dashboardID,
    cells: getCells(state, dashboardID),
  }
}

const connector = connect(mstp, null)

type Props = ConnectedProps<typeof connector>

const Cells: FunctionComponent<Props> = ({cells}) => {
  const dispatch = useDispatch()

  const onLayoutChange = (layouts: Layout[]) => {
    dispatch(updateLayout(layouts))
  }

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

export default connector(Cells)
