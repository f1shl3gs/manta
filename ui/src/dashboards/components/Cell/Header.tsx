import React, {FunctionComponent} from 'react'

interface Props {
  name: string
  note: string
  children: JSX.Element | JSX.Element[]
}

const CellHeader: FunctionComponent<Props> = ({name, note, children}) => {
  return (
    <div className="cell--header">
      <div className="cell--draggable" data-testid={`cell--draggable ${name}`}>
        <div className="cell--dot-grid" />
        <div className="cell--dot-grid" />
        <div className="cell--dot-grid" />
      </div>
      <div className="cell--name">{name}</div>
      {/*note && <CellHeaderNote note={note} />}*/}
      {note}
      {children}
    </div>
  )
}

export default CellHeader
