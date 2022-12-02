// Libraries
import React, {FunctionComponent, useState} from 'react'
import classnames from 'classnames'

// Components
import ViewOptions from 'src/visualization/TimeMachine/ViewOptions'
import {
  DraggableResizer,
  DraggableResizerPanel,
  Orientation,
} from '@influxdata/clockface'
import TimeMachineQueries from 'src/visualization/TimeMachine/TimeMachineQueries'

// Hooks
import {ViewPropertiesProvider} from 'src/visualization/TimeMachine/useViewProperties'

// Types
import {ViewProperties} from 'src/types/dashboard'
import TimeMachineVis from 'src/visualization/TimeMachine/TimeMachineVis'
import {useViewOption} from 'src/shared/useViewOption'

const INITIAL_RESIZER_HANDLE = 0.5

interface Props {
  viewProperties: ViewProperties
  onChange?: (ViewProperties) => void
}

const TimeMachine: FunctionComponent<Props> = ({viewProperties}) => {
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])
  const {isViewingVisOptions} = useViewOption()
  const containerClassName = classnames('time-machine', {
    'time-machine--split': isViewingVisOptions,
  })

  return (
    <ViewPropertiesProvider viewProperties={viewProperties}>
      {isViewingVisOptions && <ViewOptions />}

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
    </ViewPropertiesProvider>
  )
}

export default TimeMachine
