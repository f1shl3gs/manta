// Libraries
import React, {FunctionComponent, lazy} from 'react'
import {Routes, Route} from 'react-router-dom'

// Components
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

      <ChecksControlBar />

      <GetResources resources={[ResourceType.Checks]}>
        <CheckCards />
      </GetResources>
    </>
  )
}

export default ChecksIndex
