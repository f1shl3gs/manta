// Libraries
import React, { useState } from 'react';
import classnames from 'classnames';

// Components
import { DraggableResizer, Orientation } from '@influxdata/clockface';
import TimeMachineVis from './TimeMachineVis';
import TimeMachineQueries from './TimeMachineQueries';
import ViewOptions from './ViewOptions';
import { View } from '../../types';

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

  return (
    <>
      {isViewingVisOptions && <ViewOptions view={view.properties} />}

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
                {bottomContents}
              </div>
            </div>
          </DraggableResizer.Panel>
        </DraggableResizer>
      </div>
    </>
  );
};

export default TimeMachine;
