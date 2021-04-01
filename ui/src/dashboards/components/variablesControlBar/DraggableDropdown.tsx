// Libraries
import React from 'react'
import classnames from 'classnames'

// Components
import {
  DragSource,
  DropTarget,
  ConnectDropTarget,
  ConnectDragSource,
  ConnectDragPreview,
  DropTargetConnector,
  DragSourceConnector,
  DragSourceMonitor,
} from 'react-dnd'
import VariableDropdown from './VariableDropdown'

// Constants
const dropdownType = 'dropdown'

const dropdownSource = {
  beginDrag(props: Props) {
    return {
      id: props.id,
      index: props.index,
    }
  },
}

const dropdownTarget = {
  // @ts-ignore
  hover(props, monitor, component) {
    // component is always null if the draggable component is function component

    const dragIndex = monitor.getItem().index
    const hoverIndex = props.index

    // Don't replace items with themselves
    if (dragIndex === hoverIndex) {
      return
    }

    // Time to actually perform the action
    props.moveDropdown(dragIndex, hoverIndex)

    monitor.getItem().index = hoverIndex
  },
}

interface Props {
  id: string
  index: number
  name: string
  moveDropdown: (dragIndex: number, hoverIndex: number) => void
}

interface DropdownSourceCollectedProps {
  isDragging: boolean
  connectDragSource: ConnectDragSource
  connectDragPreview: ConnectDragPreview
}

interface DropdownTargetCollectedProps {
  connectDropTarget?: ConnectDropTarget
}

const DraggableDropdown: React.FC<
  Props & DropdownSourceCollectedProps & DropdownTargetCollectedProps
> = props => {
  const {
    id,
    name,
    isDragging,
    connectDragSource,
    connectDropTarget,
    connectDragPreview,
  } = props

  const className = classnames('variable-dropdown', {
    'variable-dropdown__dragging': isDragging,
  })

  // @ts-ignore
  return connectDropTarget(
    <div className={'variable-dropdown--container'}>
      {connectDragPreview(
        <div className={className}>
          {/*  */}
          <div className={'variable-dropdown--label'}>
            {connectDragSource(
              <div className={'variable-dropdown--drag'}>
                <span className={'hamburger'} />
              </div>
            )}
            <span>{name}</span>
          </div>
          <VariableDropdown variableID={id} />
        </div>
      )}
      <div className={'variable-dropdown--placeholder'} />
    </div>
  )
}

export default DropTarget<Props, DropdownTargetCollectedProps>(
  dropdownType,
  dropdownTarget,
  (connect: DropTargetConnector) => ({
    connectDropTarget: connect.dropTarget(),
  })
)(
  DragSource<Props, DropdownSourceCollectedProps>(
    dropdownType,
    dropdownSource,
    (connect: DragSourceConnector, monitor: DragSourceMonitor) => ({
      connectDragSource: connect.dragSource(),
      connectDragPreview: connect.dragPreview(),
      isDragging: monitor.isDragging(),
    })
  )(DraggableDropdown)
)
