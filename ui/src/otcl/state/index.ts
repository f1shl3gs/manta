import { useCallback, useState } from 'react';
import constate from 'constate';
import useFetch, { CachePolicies } from 'shared/hooks/useFetch';

import { Otcl } from 'types';
import remoteDataState from 'utils/rds';

const OtclPrefix = '/api/v1/otcls';

const emptyOtcl: Otcl = {
  id: '',
  name: '',
  desc: '',
  content: '',
  created: '',
  modified: '',
};

type State = {
  orgID: string;
};

const [OtclProvider, useOtcls, useOtcl] = constate(
  ({ orgID }: State) => {
    const [mutated, setMutated] = useState<number>(0);
    const [otcl, setOtcl] = useState<Otcl>(emptyOtcl);
    const reload = useCallback(() => {
      setMutated((mutated) => mutated + 1);
    }, []);

    return {
      otcl,
      setOtcl,
      orgID,
      reload,
      mutated,
    };
  },
  // useOtcls
  (value) => {
    const { orgID, mutated } = value;

    const { data, error, loading } = useFetch(
      `${OtclPrefix}?orgID=${orgID}`,
      {
        cachePolicy: CachePolicies.NO_CACHE,
      },
      [mutated]
    );

    return {
      reload: value.reload,
      otcls: data,
      rds: remoteDataState(loading, error),
    };
  },
  // useOtcl
  (value) => {
    return {
      otcl: value.otcl,
      setOtcl: value.setOtcl,
    };
  }
);

export { OtclProvider, useOtcls, useOtcl, emptyOtcl };
