import React, { useCallback } from 'react';
import ViewEditorOverlayHeader from './ViewEditorOverlayHeader';
import TimeMachine from '../../components/timeMachine/TimeMachine';
import { useCell } from './useCell';
import { useViewProperties } from '../../shared/useViewProperties';
import { Cell } from '../../types/Dashboard';

const ViewEditor: React.FC = () => {
  console.log('render')
  const { viewProperties } = useViewProperties();
  const { cell, updateCell } = useCell();
  const onNameSet = (name: string) => {
    console.log('on name set', name);
  };

  const onCancel = () => {

  };

  const onSave = useCallback(() => updateCell({
    ...cell,
    viewProperties
  } as Cell), [cell, viewProperties]);

  return (
    <>
      <ViewEditorOverlayHeader
        name={cell?.name || ''}
        onNameSet={onNameSet}
        onSave={onSave}
        onCancel={onCancel}
      />

      <div className={'veo-contents'}>
        <TimeMachine viewProperties={viewProperties} />
      </div>
    </>
  );
};

export default ViewEditor;
