import React from 'react'

import CellHeader from './CellHeader'
import CellContext from './CellContext'

import {Cell} from 'types'

interface Props {
  cell: Cell
}

const CellComponent: React.FC<Props> = ({cell}) => {
  return (
    <>
      <CellHeader name={cell.name || 'Name this Cell'} note={''}>
        <CellContext cell={cell} view={cell.properties} />
      </CellHeader>

      <div className="cell--view">
        <div>{cell.id}</div>
      </div>
    </>
  )
}

export default CellComponent
