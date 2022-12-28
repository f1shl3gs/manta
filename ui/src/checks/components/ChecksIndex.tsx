// Libraries
import React, {FunctionComponent, lazy} from 'react'
import {Routes, Route} from 'react-router-dom'

// Components
import {Page, PageContents, PageHeader, PageTitle} from '@influxdata/clockface'
import GetResources from 'src/resources/components/GetResources'
import CheckCards from 'src/checks/components/CheckCards'

// Types
import {ResourceType} from 'src/types/resources'
import ChecksControlBar from './ChecksControlBar'

// Lazy load
const NewCheckOverlay = lazy(
  () => import('src/checks/components/NewCheckOverlay')
)
const EditCheckOverlay = lazy(
  () => import('src/checks/components/EditCheckOverlay')
)

const ChecksIndex: FunctionComponent = () => {
  return (
    <>
      <Routes>
        <Route path="new" element={<NewCheckOverlay />} />
        <Route path=":id/edit" element={<EditCheckOverlay />} />
      </Routes>

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
    </>
  )
}

export default ChecksIndex
