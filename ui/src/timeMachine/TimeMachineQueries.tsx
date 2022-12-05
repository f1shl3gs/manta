// Libraries
import React, {FunctionComponent, useCallback} from 'react'
import {useDispatch, useSelector} from 'react-redux'

// Components
import QueryTabs from 'src/timeMachine/QueryTabs'
import {FlexBox} from '@influxdata/clockface'
import QueryEditor from 'src/timeMachine/QueryEditor'
import SubmitQueryButton from 'src/timeMachine/SubmitQueryButton'
import {AppState} from 'src/types/stores'
import {setActiveQueryText} from 'src/timeMachine/actions'

const TimeMachineQueries: FunctionComponent = () => {
  const dispatch = useDispatch()
  const activeQuery = useSelector((state: AppState) => {
    const {viewProperties, activeQueryIndex} = state.timeMachine
    return viewProperties.queries[activeQueryIndex]
  })
  const handleOnChange = useCallback(
    (text: string) => {
      dispatch(setActiveQueryText(text))
    },
    [dispatch]
  )

  return (
    <div className={'time-machine-queries'}>
      <div className={'time-machine-queries--controls'}>
        <QueryTabs />

        <FlexBox>
          <SubmitQueryButton />
        </FlexBox>
      </div>

      <div className={'time-machine-queries--body'}>
        <QueryEditor query={activeQuery.text || ''} onChange={handleOnChange} />
      </div>
    </div>
  )
}

export default TimeMachineQueries
