import React, {useState} from 'react'
import {useCheck} from './useCheck'
import {DraggableResizer, Orientation} from '@influxdata/clockface'
import classnames from 'classnames'
import {TimeRangeProvider} from '../shared/useTimeRange'
import CheckVis from './CheckVis'
import {AutoRefreshProvider} from '../shared/useAutoRefresh'

const INITIAL_RESIZER_HANDLE = 0.65

const CheckEditor: React.FC = () => {
  const {check, remoteDataState} = useCheck()
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])
  const containerClassName = classnames('time-machine', {
    'time-machine--split': false,
  })

  return (
    <div className={'veo-contents'}>
      <div className={containerClassName}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizer.Panel>
            <div className={'time-machine--top'}>
              <TimeRangeProvider>
                <AutoRefreshProvider>
                  <CheckVis query={check.expr} />
                </AutoRefreshProvider>
              </TimeRangeProvider>
            </div>
          </DraggableResizer.Panel>

          <DraggableResizer.Panel>
            <div className={'time-machine--bottom'}>
              <div className={'time-machine--bottom-contents'}>todo</div>
            </div>
          </DraggableResizer.Panel>
        </DraggableResizer>
      </div>
    </div>
  )
}

export default CheckEditor
