// Libraries
import React from 'react'

// Components
import {
  AlignItems,
  ComponentSize,
  FlexBox,
  FlexDirection,
  Page,
  Panel,
} from '@influxdata/clockface'
import withProvider from '../utils/withProvider'
import ProfilePanelHeader from './ProfilePanelHeader'
import ProfilePanelBody from './ProfilePanelBody'
import TimelineChart from './TimelineChart'

// Hooks
import {ProfileProvider} from './useProfile'
import {ViewTypeProvider} from './useViewType'

const Title = 'Profile'

const ProfilePage: React.FC = () => {
  return (
    <Page titleTag={Title}>
      <Page.Header fullWidth={true}>
        <Page.Title title={'Profile'} testID={'profile-page--header'} />
      </Page.Header>

      <Page.Contents fullWidth={true} scrollable={true}>
        <FlexBox
          direction={FlexDirection.Column}
          margin={ComponentSize.Small}
          alignItems={AlignItems.Stretch}
          stretchToFitWidth={true}
          testID={'profile--flexbox'}
        >
          <Panel>
            <Panel.Body>
              <TimelineChart />
            </Panel.Body>
          </Panel>

          <Panel>
            <Panel.Header>
              <ProfilePanelHeader />
            </Panel.Header>

            <Panel.Body direction={FlexDirection.Row}>
              <ProfilePanelBody />
            </Panel.Body>
          </Panel>
        </FlexBox>
      </Page.Contents>
    </Page>
  )
}

export default withProvider(
  ViewTypeProvider,
  withProvider(ProfileProvider, ProfilePage)
)
