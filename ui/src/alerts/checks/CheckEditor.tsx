// Libraries
import React, {useState} from 'react'
import classnames from 'classnames'

// Components
import {
  Button,
  ComponentColor,
  DraggableResizer,
  FlexBox,
  Orientation,
} from '@influxdata/clockface'
import {TimeRangeProvider} from '../../shared/useTimeRange'
import CheckVis from './CheckVis'
import {AutoRefreshProvider} from '../../shared/useAutoRefresh'
import CheckBuilder from '../builder/CheckBuilder'

// Hooks
import {useCheck} from './useCheck'

// Constants
import {INITIAL_RESIZER_HANDLE} from '../../constants/timeMachine'

const CheckEditor: React.FC = () => {
  const {expr} = useCheck()
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])

  return (
    <div className={'veo-contents'}>
      <div className={'time-machine'}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizer.Panel>
            <div className={'time-machine--top'}>
              <TimeRangeProvider>
                <AutoRefreshProvider>
                  <CheckVis query={expr} />
                </AutoRefreshProvider>
              </TimeRangeProvider>
            </div>
          </DraggableResizer.Panel>

          <DraggableResizer.Panel>
            <div className={'time-machine--bottom'}>
              <div className={'time-machine--bottom-contents'}>
                <div className={'time-machine-queries'}>
                  <div className={'time-machine-queries--controls'}>
                    <div className={'time-machine-queries--tabs'}>Todo</div>

                    <FlexBox>
                      <Button
                        text={'submit'}
                        color={ComponentColor.Primary}
                        onClick={() => console.log('submit')}
                      />
                    </FlexBox>
                  </div>

                  <div className={'time-machine-queries--body'}>
                    <CheckBuilder />
                  </div>
                </div>
              </div>
            </div>
          </DraggableResizer.Panel>
        </DraggableResizer>
      </div>
    </div>
  )
}

export default CheckEditor
