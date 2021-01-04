// Libraries
import React from 'react';
import { Overlay, RemoteDataState, SpinnerContainer, TechnoSpinner } from '@influxdata/clockface';
import ViewEditorOverlayHeader from './ViewEditorOverlayHeader';
import { useFetch } from 'use-http';
import { useHistory, useParams } from 'react-router-dom';


import { Cell, ViewProperties, XYViewProperties } from 'types/Dashboard';

import remoteDataState from 'utils/rds';
import { CellProvider, useCell } from './useCell';
import ViewEditor from './ViewEditor';
import { ViewPropertiesProvider } from 'shared/useViewProperties';
import { ViewOptionProvider } from '../../shared/useViewOption';

interface Props {
}

const ViewEditorOverlay: React.FC<Props> = (props) => {
  const history = useHistory();
  const { cell, remoteDataState } = useCell();

  const onNameSet = (name: string) => {
    console.log('onNameSet', name);
  };

  const onCancel = () => {
    history.goBack();
  };

  const { dashboardID, cellID } = useParams<{
    cellID: string
    dashboardID: string
  }>();
  const view = {
    cellID,
    dashboardID,
    name: cell?.name as string,
    properties: cell?.viewProperties || {
      type: 'xy',
      xColumn: 'time',
      yColumn: 'value',
      axes: {
        x: {},
        y: {}
      },
      queries: [
        {
          text: 'aaa'
        }
      ]
    }
  };

  if (!cell && remoteDataState === RemoteDataState.Done) {

  }

  console.log('rmotedatastate', remoteDataState, cell)

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <SpinnerContainer spinnerComponent={<TechnoSpinner />} loading={remoteDataState}>
          <ViewPropertiesProvider viewProperties={cell?.viewProperties as ViewProperties} >
            <ViewEditor />
          </ViewPropertiesProvider>
        </SpinnerContainer>
      </div>
    </Overlay>
  );
};

const wrapper = () => (
  <ViewOptionProvider>
    <CellProvider>
      <ViewEditorOverlay />
    </CellProvider>
  </ViewOptionProvider>
);

export default wrapper;
