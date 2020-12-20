import React, { useState } from "react";

import { ErrorHandling } from "shared/decorators/errors";
import classnames from "classnames";
import { DraggableResizer, Orientation } from "@influxdata/clockface";
import TimeMachineVis from "./TimeMachineVis";
import TimeMachineQueries from "./TimeMachineQueries";

const INITIAL_RESIZER_HANDLE = 0.5;

interface Props {
  isViewingVisOptions: boolean
}

const ViewOptions = () => {
  return (
    <div>ViewOptions</div>
  );
};


const TimeMachine: React.FC<Props> = props => {
  const {
    isViewingVisOptions
  } = props;

  const [dragPosition, setDragPosition] = useState([INITIAL_RESIZER_HANDLE]);

  const containerClassName = classnames("time-machine", {
    "time-machine--split": isViewingVisOptions
  });

  const bottomContents = <TimeMachineQueries />;

  return (
    <>
      {/*
      {isViewingVisOptions && <ViewOptions />}
      */}
      <div className={containerClassName}>
        <DraggableResizer
          handleOrientation={Orientation.Horizontal}
          handlePositions={dragPosition}
          onChangePositions={setDragPosition}
        >
          <DraggableResizer.Panel>
            <div className={"time-machine--top"}>
              <TimeMachineVis />
            </div>
          </DraggableResizer.Panel>

          <DraggableResizer.Panel>
            <div
              className={"time-machine--bottom"}
            >
              <div className={"time-machine--bottom-contents"}>
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
