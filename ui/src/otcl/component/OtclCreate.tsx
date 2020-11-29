import React, { useCallback, useEffect } from 'react';
import OtclForm from './OtclForm';
import OtclOverlay from './OtclOverlay';
import { useHistory } from 'react-router';
import { useOtcl, emptyOtcl } from '../state';

type Props = {
  onSubmit: () => void;
};

const OtclCreate: React.FC<Props> = (props) => {
  const { onSubmit } = props;
  const history = useHistory();
  const { setOtcl } = useOtcl();
  const onDismiss = useCallback(() => history.goBack(), []);

  useEffect(() => {
    return () => {
      setOtcl(emptyOtcl);
    };
  });

  return (
    <OtclOverlay title={'Create new Otcl'} onDismiss={onDismiss}>
      <OtclForm onDismiss={onDismiss} onSubmit={onSubmit} />
    </OtclOverlay>
  );
};

export default OtclCreate;
