// Libraries
import React, {FunctionComponent} from 'react'

// Components
import QueryTabs from 'src/visualization/TimeMachine/QueryTabs'
import {FlexBox} from '@influxdata/clockface'
import QueryEditor from 'src/visualization/TimeMachine/QueryEditor'
import SubmitQueryButton from 'src/visualization/TimeMachine/SubmitQueryButton'

// Hooks
import {useQueries} from 'src/visualization/TimeMachine/useQueries'

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
