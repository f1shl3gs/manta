// Libraries
import React, {useState} from 'react'

// Components
import {
  AlignItems,
  Columns,
  ComponentSize,
  FlexBox,
  FlexDirection,
  Grid,
  Page,
  Panel,
} from '@influxdata/clockface'
import FlameGraphControlBar from './FlameGraphControlBar'
import TablePanel from './components/TablePanel'
import FlameGraphPanel from './components/FlameGraphPanel'
import withProvider from '../utils/withProvider'
import {ProfileProvider} from './useProfile'
import ProfilePanelHeader from './ProfilePanelHeader'

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
            <Panel.Body>Timeline</Panel.Body>
          </Panel>

          <Panel>
            <Panel.Header>
              <ProfilePanelHeader />
            </Panel.Header>

            <Panel.Body direction={FlexDirection.Row}>
              <TablePanel />
              <FlameGraphPanel />
            </Panel.Body>
          </Panel>
        </FlexBox>
      </Page.Contents>
    </Page>
  )
}

export default withProvider(ProfileProvider, ProfilePage)
