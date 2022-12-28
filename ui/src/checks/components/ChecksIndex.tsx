// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Page, PageContents, PageHeader, PageTitle} from '@influxdata/clockface'
import GetResources from 'src/resources/components/GetResources'
import CheckCards from 'src/checks/components/CheckCards'

// Types
import {ResourceType} from 'src/types/resources'
import ChecksControlBar from './ChecksControlBar'

const ChecksIndex: FunctionComponent = () => {
  return (
    <Page titleTag={'Checks'}>
      <PageHeader fullWidth={false}>
        <PageTitle title={'Checks'} />
      </PageHeader>

      <ChecksControlBar />

      <PageContents>
        <GetResources resources={[ResourceType.Checks]}>
          <CheckCards />
        </GetResources>
      </PageContents>
    </Page>
  )
}

export default ChecksIndex
