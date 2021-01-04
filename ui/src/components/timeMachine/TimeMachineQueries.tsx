// Libraries
import React, { useState } from 'react';

// Components
// import { Controlled as ReactCodeMirror } from 'react-codemirror2';
import QueryTabs from './QueryTabs';
import { FlexBox } from '@influxdata/clockface';
import SubmitQueryButton from './SubmitQueryButton';
import QueryEditor from './QueryEditor';
import { QueriesProvider, useQueries } from './useQueries';

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

  return (
    <QueriesProvider>
      <div className={'time-machine-queries'}>
        <div className={'time-machine-queries--controls'}>
          <QueryTabs />

          <FlexBox>
            <SubmitQueryButton />
          </FlexBox>
        </div>

        <div className={'time-machine-queries--body'}>
          <QueryEditor />
        </div>
      </div>
    </QueriesProvider>
  );
};

export default TimeMachineQueries;
