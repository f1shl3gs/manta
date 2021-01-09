import React, { useMemo, ReactElement } from 'react';
import FetchContext from './FetchContext';
import { FetchContextTypes, FetchProviderProps } from './types';

export const Provider = ({
  url,
  options,
  graphql = false,
  children
}: FetchProviderProps): ReactElement => {
  const defaults = useMemo(
    (): FetchContextTypes => ({
      url: url || '',
      options: options || {},
      graphql // TODO: this will make it so useFetch(QUERY || MUTATION) will work
    }),
    [options, graphql, url]
  );

  return (
    <FetchContext.Provider value={defaults}>{children}</FetchContext.Provider>
  );
};