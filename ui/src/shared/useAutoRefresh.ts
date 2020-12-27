import constate from 'constate';
import { useState } from 'react';
import { AutoRefresh, AutoRefreshStatus } from 'types/AutoRefresh';

const [AutoRefreshProvider, useAutoRefresh] = constate(
  () => {
    const [autoRefresh, setAutoRefresh] = useState<AutoRefresh>({
      status: AutoRefreshStatus.Active,
      interval: 15
    });

    return {
      autoRefresh,
      setAutoRefresh
    };
  },
  value => value
);

export {
  AutoRefreshProvider,
  useAutoRefresh
};