// Libraries
import React from 'react';
import { Overlay, SpinnerContainer, TechnoSpinner } from '@influxdata/clockface';
import ViewEditorOverlayHeader from './ViewEditorOverlayHeader';
import { useFetch } from 'use-http';
import { useHistory, useParams } from 'react-router-dom';

import TimeMachine from 'components/timeMachine/TimeMachine';
import { useViewEditor, ViewEditorProvider } from './useViewEditor';

import { Cell, XYViewProperties } from 'types/Dashboard';

import remoteDataState from 'utils/rds';

interface Props {
}

const ViewEditorOverlay: React.FC<Props> = (props) => {
  const { dashboardID, cellID } = useParams<{
    cellID: string
    dashboardID: string
  }>();
  const { data, error, loading } = useFetch<Cell>(
    `/api/v1/dashboards/${dashboardID}/cells/${cellID}`,
    {},
    []
  );
  const rds = remoteDataState(loading, error);
  const cell = data;
  const history = useHistory();

  const onNameSet = (name: string) => {
    console.log('onNameSet', name);
  };

  const onSave = () => {
    console.log('onSave');
    history.goBack();
  };

  const onCancel = () => {
    history.goBack();
  };

  const { isViewingVisOptions } = useViewEditor();
  const view = {
    cellID,
    dashboardID,
    name: cell?.name as string,
    properties: {
      type: 'xy'
    } as XYViewProperties
  };

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <SpinnerContainer spinnerComponent={<TechnoSpinner />} loading={rds}>

          <ViewEditorOverlayHeader
            name={cell?.name || ''}
            onNameSet={onNameSet}
            onSave={onSave}
            onCancel={onCancel}
          />

          <div className={'veo-contents'}>
            <TimeMachine isViewingVisOptions={isViewingVisOptions} view={view} />
          </div>
        </SpinnerContainer>
      </div>
    </Overlay>
  );
};

const wrapper = () => (
  <ViewEditorProvider>
    <ViewEditorOverlay />
  </ViewEditorProvider>
);

export default wrapper;
