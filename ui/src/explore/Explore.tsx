// Libraries
import React, {FunctionComponent, useEffect} from 'react'

// Components
import {
  ComponentSize,
  FlexBox,
  Page,
  PageContents,
  PageHeader,
  PageTitle,
} from '@influxdata/clockface'
import TimeMachine from 'src/timeMachine/components/TimeMachine'
import TimeRangeDropdown from 'src/shared/components/TimeRangeDropdown'

// Types
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'
import {useDispatch} from 'react-redux';

// Actions
import {resetTimeMachine} from 'src/timeMachine/actions'

const Explore: FunctionComponent = () => {
  const dispatch = useDispatch()

  useEffect(() => {
    return () => {
      dispatch(resetTimeMachine())
    }
  }, [dispatch])

  return (
    <Page titleTag={'Explore'}>
      <PageHeader fullWidth={true}>
        <PageTitle title={'Explore'} />

        <FlexBox margin={ComponentSize.Small}>
          <TimeRangeDropdown />
          <AutoRefreshButton />
        </FlexBox>
      </PageHeader>

      <PageContents>
        <div className={'explore-contents'}>
          <TimeMachine />
        </div>
      </PageContents>
    </Page>
  )
}

export default Explore
