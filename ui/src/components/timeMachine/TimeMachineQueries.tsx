// Libraries
import React from 'react'

// Components
import {FlexBox} from '@influxdata/clockface'
import QueryTabs from './QueryTabs'
import SubmitQueryButton from './SubmitQueryButton'
import QueryEditor from './QueryEditor'

// Hooks
import {useActiveQuery} from './useQueries'

const TimeMachineQueries: React.FC = () => {
  const {activeQuery, onSetText} = useActiveQuery()

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
