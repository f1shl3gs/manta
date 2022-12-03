// Libraries
import React, {FunctionComponent} from 'react'

// Components
import QueryTabs from 'src/timeMachine/QueryTabs'
import {FlexBox} from '@influxdata/clockface'
import QueryEditor from 'src/timeMachine/QueryEditor'
import SubmitQueryButton from 'src/timeMachine/SubmitQueryButton'

// Hooks
import {useQueries} from 'src/timeMachine/useTimeMachine'

const TimeMachineQueries: FunctionComponent = () => {
  const {onSetText, activeQuery} = useQueries()

  return (
    <div className={'time-machine-queries'}>
      <div className={'time-machine-queries--controls'}>
        <QueryTabs />

        <FlexBox>
          <SubmitQueryButton />
        </FlexBox>
      </div>

      <div className={'time-machine-queries--body'}>
        <QueryEditor query={activeQuery.text || ''} onChange={onSetText} />
      </div>
    </div>
  )
}

export default TimeMachineQueries
