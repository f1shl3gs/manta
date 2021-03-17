// Libraries
import React, {useState} from 'react'

import {useCheck} from './useCheck'

// Components
import {
  Button,
  ComponentColor,
  DraggableResizer,
  FlexBox,
  Orientation,
} from '@influxdata/clockface'
import {Controlled as ReactCodeMirror} from 'react-codemirror2'

import classnames from 'classnames'
import {TimeRangeProvider} from '../../shared/useTimeRange'
import CheckVis from './CheckVis'
import {AutoRefreshProvider} from '../../shared/useAutoRefresh'
import CheckBuilder from '../builder/CheckBuilder'

const INITIAL_RESIZER_HANDLE = 0.65

const options = {
  tabIndex: 1,
  mode: 'yaml',
  readonly: true,
  lineNumbers: true,
  autoRefresh: true,
  theme: 'material',
  completeSingle: false,
}

const CheckEditor: React.FC = () => {
  const {check} = useCheck()
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

          <DraggableResizer.Panel sizePercent={30}>
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
