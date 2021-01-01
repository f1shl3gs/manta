// Libraries
import React, { useState } from 'react';

// Components
import { Controlled as ReactCodeMirror } from 'react-codemirror2';

// Constants
const options = {
  tabIndex: 1,
  mode: '',
  lineNumber: true,
  autoRefresh: true,
  theme: 'time-machine',
  completeSingle: false
};

const TimeMachineQueries: React.FC = () => {
  const [content, setContent] = useState('');

  return (
    <div style={{ height: '100%' }}>
      <ReactCodeMirror
        autoCursor={true}
        autoScroll={true}
        value={content}
        options={options}
        onBeforeChange={(editor, data, value) => setContent(value)}
      />
    </div>
  );
};

export default TimeMachineQueries;
