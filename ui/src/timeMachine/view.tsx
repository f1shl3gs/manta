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
import {useSelector} from 'react-redux'

// Types
import TimeMachineVis from 'src/timeMachine/TimeMachineVis'
import {AppState} from 'src/types/stores'

const INITIAL_RESIZER_HANDLE = 0.5

const TimeMachine: FunctionComponent = () => {
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])
  const viewingVisOptions = useSelector((state: AppState) => {
    return state.timeMachine.viewingVisOptions
  })
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
              <TimeMachineVis />
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
