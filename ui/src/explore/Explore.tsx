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
import TimeMachine from 'src/visualization/TimeMachine'
import TimeRangeDropdown from 'src/dashboards/components/TimeRangeDropdown'

// Types
import AutoRefreshButton from 'src/shared/components/AutoRefreshButton'

// Constants
import {defaultViewProperties} from 'src/constants/dashboard'

const Explore: FunctionComponent = () => {
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
          <TimeMachine viewProperties={defaultViewProperties} />
        </div>
      </PageContents>
    </Page>
  )
}

export default Explore
