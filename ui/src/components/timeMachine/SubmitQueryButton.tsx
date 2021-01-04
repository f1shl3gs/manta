import React from 'react';
import { Button, ComponentColor } from '@influxdata/clockface';

const SubmitQueryButton: React.FC = () => {
  return (
    <Button
      text={'Submit'}
      color={ComponentColor.Primary}
      onClick={() => console.log('submit')}
    />
  );
};

export default SubmitQueryButton;
