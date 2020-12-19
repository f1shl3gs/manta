import React from "react";
import { Cell } from "../../types";
import CellHeader from "./CellHeader";
import CellContext from "./CellContext";

// Style
import './CellComponent.scss';

interface Props {
  cell: Cell
}

const CellComponent: React.FC<Props> = ({ cell }) => {
  return (
    <>
      <CellHeader name={cell.name} note={""}>
        <CellContext
          cell={cell}
          view={cell.properties}
        />
      </CellHeader>

      <div className="cell--view">
        <div>{cell.id}</div>
      </div>
    </>
  );
};

export default CellComponent;
