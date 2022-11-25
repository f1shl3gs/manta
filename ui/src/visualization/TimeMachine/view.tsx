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
import useToggle from 'src/shared/useToggle'
import {ViewPropertiesProvider} from 'src/visualization/TimeMachine/useViewProperties'

// Types
import {ViewProperties} from 'src/types/Dashboard'
import {QueriesProvider} from 'src/visualization/TimeMachine/useQueries'
import TimeMachineVis from 'src/visualization/TimeMachine/TimeMachineVis'

const INITIAL_RESIZER_HANDLE = 0.5

interface Props {
  viewProperties: ViewProperties
  onChange?: (ViewProperties) => void
}

const TimeMachine: FunctionComponent<Props> = ({viewProperties}) => {
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])
  const [showVisOptions] = useToggle(false)
  const containerClassName = classnames('time-machine', {
    'time-machine--split': showVisOptions,
  })

  return (
    <ViewPropertiesProvider viewProperties={viewProperties}>
      {showVisOptions && <ViewOptions />}

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
                <QueriesProvider>
                  <TimeMachineQueries />
                </QueriesProvider>
              </div>
            </div>
          </DraggableResizerPanel>
        </DraggableResizer>
      </div>
    </ViewPropertiesProvider>
  )
}

export default TimeMachine
