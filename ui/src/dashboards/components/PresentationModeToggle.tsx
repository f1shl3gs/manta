import React from 'react';
import { IconFont, SquareButton } from '@influxdata/clockface';
import { usePresentationMode } from '../../shared/usePresentationMode';

const PresentationModeToggle = () => {
  const { togglePresentationMode } = usePresentationMode();

  return (
    <SquareButton
      icon={IconFont.ExpandA}
      onClick={togglePresentationMode}
    />
  );
};

export default PresentationModeToggle;