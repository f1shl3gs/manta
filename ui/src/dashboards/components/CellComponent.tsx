import React from 'react'

import CellHeader from './CellHeader'
import CellContext from './CellContext'

import {Cell} from 'types/Dashboard'
import TimeSeries from './TimeSeries';

interface Props {
  cell: Cell
}

const CellComponent: React.FC<Props> = ({cell}) => {
  return (
    <>
      <CellHeader name={cell.name || 'Name this Cell'} note={''}>
        <CellContext cell={cell} view={cell.viewProperties} />
      </CellHeader>

      <div className="cell--view">
        <TimeSeries/>
      </div>
    </>
  )
}

export default CellComponent
