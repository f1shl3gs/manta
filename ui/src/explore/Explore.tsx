// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ComponentSize,
  FlexBox,
  Page,
  PageContents,
  PageHeader,
  PageTitle,
} from '@influxdata/clockface'
import TimeMachine from 'src/timeMachine'
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'
import TimeRangeDropdown from 'src/dashboards/components/TimeRangeDropdown'
import {TimeMachineProvider} from 'src/timeMachine/useTimeMachine'

const Explore: FunctionComponent = () => {
  return (
    <Page titleTag={'Explore'}>
      <TimeMachineProvider>
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
      </TimeMachineProvider>
    </Page>
  )
}

export default Explore
