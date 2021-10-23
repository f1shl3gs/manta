// Libraries
import React from 'react'

// Components
import {FlexBox} from '@influxdata/clockface'
import QueryTabs from './QueryTabs'
import SubmitQueryButton from './SubmitQueryButton'
import QueryEditor from './QueryEditor'

// Hooks
import {useQueries} from './useQueries'

const TimeMachineQueries: React.FC = () => {
  const {activeIndex, queries, onSetText} = useQueries()
  console.log('queries', queries, 'active', activeIndex)
  const activeQuery = queries[activeIndex]

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
