import React from 'react';
import { useActiveQuery } from './useQueries';

const QueryTabName: React.FC = () => {
  const { activeQuery } = useActiveQuery();
  const { name = 'Query' } = activeQuery;

  return (
    <div className={'query-tab--name'} title={name}>
      {name}
    </div>
  );
};

export default QueryTabName;