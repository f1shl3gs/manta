import React, { useCallback, useEffect } from 'react';
import { useHistory, useParams } from 'react-router';
import useFetch from 'use-http';

import { SpinnerContainer, TechnoSpinner } from '@influxdata/clockface';
import remoteDataState from 'utils/rds';
import { emptyOtcl, useOtcl } from 'otcl/state';

import OtclForm from './OtclForm';
import OtclOverlay from './OtclOverlay';

const OtclEdit: React.FC = () => {
  const { otclID } = useParams();
  const { otcl, setOtcl } = useOtcl();

  console.log('otcl', otcl);

  const { get, patch, loading, error } = useFetch(
    `/api/v1/otcls/${otclID}`,
    {}
  );

  const history = useHistory();
  const onDismiss = useCallback(() => history.goBack(), []);
  const onSubmit = useCallback(() => {
    console.log('submit', otcl);
    patch(otcl).finally(() => {
      console.log('submit');
    });
  }, [otcl]);

  useEffect(() => {
    get().then((resp) => {
      setOtcl(resp);
    });

    return () => {
      setOtcl(emptyOtcl);
    };
  }, [otclID]);

  const rds = remoteDataState(loading, error);

  return (
    <OtclOverlay title={'Update Otcl Config'} onDismiss={onDismiss}>
      <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
        <OtclForm onSubmit={onSubmit} onDismiss={onDismiss} />
      </SpinnerContainer>
    </OtclOverlay>
  );
};

export default OtclEdit;
