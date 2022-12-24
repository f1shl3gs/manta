// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import QueryTabs from 'src/timeMachine/components/queryEditor/QueryTabs'
import {FlexBox} from '@influxdata/clockface'
import QueryEditor from 'src/timeMachine/components/queryEditor/QueryEditor'
import SubmitQueryButton from 'src/timeMachine/components/SubmitQueryButton'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {setActiveQueryText} from 'src/timeMachine/actions'

const mstp = (state: AppState) => {
  const {viewProperties, activeQueryIndex} = state.timeMachine

  return {
    query: viewProperties.queries[activeQueryIndex].text,
  }
}

const mdtp = {
  onChange: setActiveQueryText,
}

const connector = connect(mstp, mdtp)

type Props = ConnectedProps<typeof connector>

const TimeMachineQueries: FunctionComponent<Props> = ({query, onChange}) => {
  return (
    <div className={'time-machine-queries'}>
      <div className={'time-machine-queries--controls'}>
        <QueryTabs />

        <FlexBox>
          <SubmitQueryButton />
        </FlexBox>
      </div>

      <div className={'time-machine-queries--body'}>
        <QueryEditor query={query} onChange={onChange} />
      </div>
    </div>
  )
}

export default connector(TimeMachineQueries)
