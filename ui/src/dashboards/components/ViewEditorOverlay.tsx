// Libraries
import React from "react";
import { Overlay, RemoteDataState, SpinnerContainer, TechnoSpinner } from "@influxdata/clockface";
import ViewEditorOverlayHeader from "./ViewEditorOverlayHeader";
import { useFetch } from "use-http";
import { useHistory, useParams } from "react-router-dom";
import remoteDataState from "../../utils/rds";

import { Cell } from "types";
import TimeMachine from "../../components/timeMachine/TimeMachine";

interface Props {

}

const ViewEditorOverlay: React.FC<Props> = props => {
  const { dashboardID, cellID } = useParams<{ dashboardID: string, cellID: string }>();
  const { data, error, loading } = useFetch<Cell>(`/api/v1/dashboards/${dashboardID}/cells/${cellID}`, {}, []);
  const rds = remoteDataState(loading, error);
  const cell = data;
  const history = useHistory();

  const onNameSet = (name: string) => {
    console.log("onNameSet", name);
  };

  const onSave = () => {
    console.log("onSave");
    history.goBack();
  };

  const onCancel = () => {
    history.goBack();
  };

  return (
    <Overlay visible={true} className={"veo-overlay"}>
      <div className={"veo"}>
        <SpinnerContainer
          spinnerComponent={<TechnoSpinner />}
          loading={rds}
        >
          <ViewEditorOverlayHeader
            name={cell?.name || ""}
            onNameSet={onNameSet}
            onSave={onSave}
            onCancel={onCancel}
          />

          <div className={"veo-contents"}>
            <TimeMachine isViewingVisOptions={true}/>
          </div>
        </SpinnerContainer>
      </div>
    </Overlay>
  );
};

export default ViewEditorOverlay;
