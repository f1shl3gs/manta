// Libraries
import React, { useState } from 'react';
import classnames from 'classnames';

// Components
import { DraggableResizer, Orientation } from '@influxdata/clockface';
import TimeMachineVis from './TimeMachineVis';
import TimeMachineQueries from './TimeMachineQueries';
import ViewOptions from './ViewOptions';

// Types
import { View, ViewProperties, XYViewProperties } from 'types/Dashboard';
import { ViewPropertiesProvider } from './useView';

const INITIAL_RESIZER_HANDLE = 0.5;

interface Props {
  isViewingVisOptions: boolean,
  view: View
}

const TimeMachine: React.FC<Props> = (props) => {
  const { isViewingVisOptions, view } = props;

  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE]);

  const containerClassName = classnames('time-machine', {
    'time-machine--split': isViewingVisOptions
  });

  const bottomContents = <TimeMachineQueries />;

  const vp = {
    type: 'xy',
    xColumn: 'time',
    yColumn: 'value',
    axes: {
      x: {

      },
      y: {

      }
    },
    queries: [
      {
        text: 'aaa'
      }
    ]
  } as XYViewProperties

  return (
    <ViewPropertiesProvider viewProperties={vp as ViewProperties}>
      {isViewingVisOptions && <ViewOptions/>}

      <div className={containerClassName}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizer.Panel>
            <div className={'time-machine--top'}>
              <TimeMachineVis viewProperties={vp}/>
            </div>
          </DraggableResizer.Panel>

          <DraggableResizer.Panel>
            <div className={'time-machine--bottom'}>
              <div className={'time-machine--bottom-contents'}>
                {bottomContents}
              </div>
            </div>
          </DraggableResizer.Panel>
        </DraggableResizer>
      </div>
    </ViewPropertiesProvider>
  );
};

export default TimeMachine;
