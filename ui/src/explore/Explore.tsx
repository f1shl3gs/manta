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
import {ViewProperties} from 'src/types/Dashboard'
import AutoRefreshButton from '../shared/components/AutoRefreshButton'

const defaultViewProperties: ViewProperties = {
  type: 'xy',
  xColumn: 'time',
  yColumn: 'value',
  hoverDimension: 'auto',
  geom: 'line',
  position: 'overlaid',
  axes: {
    x: {},
    y: {},
  },
  queries: [
    {
      name: 'query 1',
      text: '',
      hidden: false,
    },
  ],
}

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
