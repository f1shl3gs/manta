import React from 'react'
import {DndProvider} from 'react-dnd'
import {HTML5Backend} from 'react-dnd-html5-backend'

function withDragDropContext<TProps extends {}>(Component: React.FC<TProps>) {
  return (props: TProps) => (
    <DndProvider backend={HTML5Backend}>
      <Component {...props} />
    </DndProvider>
  )
}

export default withDragDropContext
