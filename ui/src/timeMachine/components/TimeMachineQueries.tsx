// Libraries
import React, {FunctionComponent} from 'react'

// Components
import QueryTabs from 'src/timeMachine/components/QueryTabs'
import {FlexBox} from '@influxdata/clockface'
import QueryEditor from 'src/timeMachine/components/QueryEditor'
import SubmitQueryButton from 'src/timeMachine/components/SubmitQueryButton'

const TimeMachineQueries: FunctionComponent = () => {
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
