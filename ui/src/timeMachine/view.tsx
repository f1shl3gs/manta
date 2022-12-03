// Libraries
import React, {FunctionComponent, useState} from 'react'
import classnames from 'classnames'

// Components
import ViewOptions from 'src/timeMachine/ViewOptions'
import {
  DraggableResizer,
  DraggableResizerPanel,
  Orientation,
} from '@influxdata/clockface'
import TimeMachineQueries from 'src/timeMachine/TimeMachineQueries'

// Hooks
import {useTimeMachine} from 'src/timeMachine/useTimeMachine'

// Types
import TimeMachineVis from 'src/timeMachine/TimeMachineVis'

const INITIAL_RESIZER_HANDLE = 0.5

const TimeMachine: FunctionComponent = () => {
  const {viewProperties, viewingVisOptions} = useTimeMachine()
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])
  const containerClassName = classnames('time-machine', {
    'time-machine--split': viewingVisOptions,
  })

  return (
    <>
      {viewingVisOptions && <ViewOptions />}

      <div className={containerClassName}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizerPanel>
            <div className={'time-machine--top'}>
              <TimeMachineVis viewProperties={viewProperties} />
            </div>
          </DraggableResizerPanel>

          <DraggableResizerPanel>
            <div className={'time-machine--bottom'}>
              <div className={'time-machine--bottom-contents'}>
                <TimeMachineQueries />
              </div>
            </div>
          </DraggableResizerPanel>
        </DraggableResizer>
      </div>
    </>
  )
}

export default TimeMachine
