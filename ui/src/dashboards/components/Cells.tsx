import React from "react";
import ReactGridLayout, { Layout, WidthProvider } from "react-grid-layout";
import { useCells } from "../state/dashboard";
import { Cell } from "../../types";
import CellComponent from "./CellComponent";

const Grid = WidthProvider(ReactGridLayout);

// styles

const DASHBOARD_LAYOUT_ROW_HEIGHT = 83.5;
const LAYOUT_MARGIN = 4;

const layout = (cells: Cell[]): Layout[] => {
  return cells.map(c => {
    return {
      i: c.id,
      w: c.w || 3,
      h: c.h || 3,
      x: c.x || 0,
      y: c.y || 0
    };
  });
};

const Cells: React.FC = () => {
  const cells = useCells();

  console.log('cells', cells)

  return (
    <>
      {/* ScrollDetector!?*/}
      <Grid
        cols={12}
        layout={layout(cells)}
        rowHeight={DASHBOARD_LAYOUT_ROW_HEIGHT}
        useCSSTransforms={false}
        containerPadding={[0, 0]}
        margin={[LAYOUT_MARGIN, LAYOUT_MARGIN]}
        draggableHandle={".cell-draggable"}
        onLayoutChange={next => console.log("layout change")}
      >
        {
          cells.map(cell => (
            <div
              key={cell.id}
              className="cell"
            >
              <CellComponent cell={cell} />
            </div>
          ))
        }
      </Grid>
    </>
  );
};

export default Cells;