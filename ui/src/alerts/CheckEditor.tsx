import React, {useState} from 'react'
import {useCheck} from './useCheck'
import {DraggableResizer, Orientation} from '@influxdata/clockface'
import CheckVis from './CheckVis'

const INITIAL_RESIZER_HANDLE = 0.5

const CheckEditor: React.FC = () => {
  const {check} = useCheck()
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])

  return (
    <>
      <DraggableResizer
        handleOrientation={Orientation.Horizontal}
        handlePositions={dragPosition}
        onChangePositions={setDragPosition}
      >
        <DraggableResizer.Panel>
          <CheckVis query={check.expr} />
        </DraggableResizer.Panel>

        <DraggableResizer.Panel>
          <div>Todo</div>
        </DraggableResizer.Panel>
      </DraggableResizer>
    </>
  )
}

export default CheckEditor
