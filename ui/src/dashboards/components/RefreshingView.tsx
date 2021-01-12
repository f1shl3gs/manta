import React, { useEffect } from 'react';
import TimeSeries from './TimeSeries';
import { useAutoRefresh } from '../../shared/useAutoRefresh';

const RefreshingView: React.FC = () => {
  const { autoRefresh } = useAutoRefresh();
  useEffect(() => {

  })

  return (
    <TimeSeries />
  );
};

export default RefreshingView;