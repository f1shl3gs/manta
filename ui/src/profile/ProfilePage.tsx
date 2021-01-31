import React from 'react'
import {Page} from '@influxdata/clockface'

const Title = 'Profile'

const ProfilePage: React.FC = () => {
  return (
    <Page titleTag={Title}>
      <Page.Contents fullWidth={false} scrollable={false}>
        Todo
      </Page.Contents>
    </Page>
  )
}

export default ProfilePage
