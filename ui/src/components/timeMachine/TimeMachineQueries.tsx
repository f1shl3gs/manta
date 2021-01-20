// Libraries
import React from 'react'

// Components
import {FlexBox} from '@influxdata/clockface'
import QueryTabs from './QueryTabs'
import SubmitQueryButton from './SubmitQueryButton'
import QueryEditor from './QueryEditor'

const TimeMachineQueries: React.FC = () => {
  return (
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
  )
}

export default TimeMachineQueries
