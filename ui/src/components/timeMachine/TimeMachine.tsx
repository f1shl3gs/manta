// Libraries
import React, {useState} from 'react'
import classnames from 'classnames'

// Components
import {DraggableResizer, Orientation} from '@influxdata/clockface'
import TimeMachineVis from './TimeMachineVis'
import TimeMachineQueries from './TimeMachineQueries'
import ViewOptions from './ViewOptions'

// Types
import {ViewProperties} from 'types/Dashboard'
import {useViewOption} from 'shared/useViewOption'
import {QueriesProvider} from './useQueries'

const INITIAL_RESIZER_HANDLE = 0.5

interface Props {
  viewProperties: ViewProperties
  bottomContents?: JSX.Element
}

const TimeMachine: React.FC<Props> = props => {
  const {isViewingVisOptions} = useViewOption()
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])

  const containerClassName = classnames('time-machine', {
    'time-machine--split': isViewingVisOptions,
  })

  return (
    <QueriesProvider>
      {isViewingVisOptions && <ViewOptions />}

      <div className={containerClassName}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizer.Panel>
            <div className={'time-machine--top'}>
              <TimeMachineVis />
            </div>
          </DraggableResizer.Panel>

          <DraggableResizer.Panel>
            <div className={'time-machine--bottom'}>
              <div className={'time-machine--bottom-contents'}>
                <TimeMachineQueries />
              </div>
            </div>
          </DraggableResizer.Panel>
        </DraggableResizer>
      </div>
    </QueriesProvider>
  )
}

export default TimeMachine
