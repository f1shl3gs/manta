// Libraries
import React from 'react';

// Components
import { Controlled as ReactCodeMirror } from 'react-codemirror2';
import { useActiveQuery } from './useQueries';

// Constants
const options = {
  tabIndex: 1,
  lineNumber: true,
  autoRefresh: true,
  theme: 'time-machine',
  completeSingle: false
};

const QueryEditor: React.FC = () => {
  const { activeQuery, onSetText } = useActiveQuery();

  return (
    <ReactCodeMirror
      autoScroll={true}
      value={activeQuery.text}
      options={options}
      onBeforeChange={(editor, data, value) => onSetText(value)}
    />
  );
};

export default QueryEditor;