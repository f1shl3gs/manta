import React from 'react';

import { Button, ComponentColor, IconFont } from '@influxdata/clockface';
import { useViewEditor } from './useViewEditor';

const VisOptionsButton: React.FC = () => {
  const { isViewingVisOptions, onToggleVisOptions } = useViewEditor();

  const color = isViewingVisOptions ?
    ComponentColor.Primary : ComponentColor.Default;

  return (
    <Button
      color={color}
      icon={IconFont.CogThick}
      onClick={onToggleVisOptions}
      text={'Customize'}
    />
  );
};

export default VisOptionsButton;