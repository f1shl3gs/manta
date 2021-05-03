// Libraries
import React, {useCallback, useState} from 'react'

// Components
import {
  Button,
  ComponentColor,
  DraggableResizer,
  FlexBox,
  Orientation,
} from '@influxdata/clockface'
import CheckVis from './CheckVis'
import {TimeRangeProvider} from 'shared/useTimeRange'
import {AutoRefreshProvider} from 'shared/useAutoRefresh'

// Hooks
import {useCheck} from './useCheck'

// Constants
import {INITIAL_RESIZER_HANDLE} from '../../constants/timeMachine'
import QueryBuilder from './builder/QueryBuilder'
import CheckBuilder from './builder/CheckBuilder'

const CheckEditor: React.FC = () => {
  const {expr, tab, onExprUpdate} = useCheck()
  const [query, setQuery] = useState(expr)
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])

  const onSubmit = useCallback(() => {
    setQuery(expr)
  }, [expr])

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
                  <CheckVis query={query} />
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
                        onClick={onSubmit}
                      />
                    </FlexBox>
                  </div>

                  <div className={'time-machine-queries--body'}>
                    {tab === 'query' ? (
                      <QueryBuilder
                        expr={expr}
                        onExprUpdate={onExprUpdate}
                        onSubmit={setQuery}
                      />
                    ) : (
                      <CheckBuilder />
                    )}
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
