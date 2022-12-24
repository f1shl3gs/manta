// Libraries
import React, {FunctionComponent, useState} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  Button,
  ComponentColor,
  ComponentSize,
  DraggableResizer,
  DraggableResizerPanel,
  FlexBox,
  Orientation,
} from '@influxdata/clockface'
import TimeSeries from 'src/shared/components/TimeSeries'

// Types
import {AppState} from 'src/types/stores'

// Utils
import {createView} from 'src/visualization/helper'
import PropertiesCard from './PropertiesCard'
import CheckQueryEditor from './CheckQueryEditor'
import TabSwitcher from './TabSwitcher'
import ConditionList from './ConditionList'

const INITIAL_RESIZER_HANDLE = 0.5

const mstp = (state: AppState) => {
  const {query, tab} = state.checkBuilder
  const viewProperties = createView('xy')

  return {
    tab,
    query,
    viewProperties: {
      ...viewProperties,
      queries: [
        {
          text: query,
          hidden: false,
        },
      ],
    },
  }
}

const connector = connect(mstp, null)
type Props = ConnectedProps<typeof connector>

const CheckBuilder: FunctionComponent<Props> = ({tab, viewProperties}) => {
  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE])

  return (
    <div className={'time-machine'}>
      <DraggableResizer
        handleOrientation={Orientation.Horizontal}
        handlePositions={dragPosition}
        onChangePositions={setDragPosition}
      >
        <DraggableResizerPanel>
          <div className={'time-machine--top'}>
            <TimeSeries viewProperties={viewProperties} />
          </div>
        </DraggableResizerPanel>

        <DraggableResizerPanel>
          <div className={'time-machine--bottom'}>
            <div className={'time-machine--bottom-contents'}>
              <div className={'time-machine-queries'}>
                <div className={'time-machine-queries--controls'}>
                  <TabSwitcher />

                  <FlexBox>
                    <Button
                      text={'Submit'}
                      size={ComponentSize.Small}
                      color={ComponentColor.Primary}
                      onClick={() => console.log('click')}
                    />
                  </FlexBox>
                </div>

                <div className={'time-machine-queries--body'}>
                  <div className={'flux-editor'}>
                    {tab === 'query' ? (
                      <CheckQueryEditor />
                    ) : (
                      <>
                        <PropertiesCard />
                        <ConditionList />
                      </>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </DraggableResizerPanel>
      </DraggableResizer>
    </div>
  )
}

export default connector(CheckBuilder)
