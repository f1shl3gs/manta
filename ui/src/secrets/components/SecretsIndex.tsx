// Libraries
import React, {FunctionComponent, lazy} from 'react'
import {Route, Routes} from 'react-router-dom'

// Components
import GetResources from 'src/resources/components/GetResources'
import SecretTabHeader from 'src/secrets/components/SecretTabHeader'
import SecretList from 'src/secrets/components/SecretList'

// Types
import {ResourceType} from 'src/types/resources'

// Lazy load components
const CreateSecretOverlay = lazy(
  () => import('src/secrets/components/CreateSecretOverlay')
)
const EditSecretOverlay = lazy(
  () => import('src/secrets/components/EditSecretOverlay')
)

const SecretsIndex: FunctionComponent = () => {
  return (
    <>
      <Routes>
        <Route path="new" element={<CreateSecretOverlay />} />
        <Route path=":key/edit" element={<EditSecretOverlay />} />
      </Routes>

      <SecretTabHeader />

      <GetResources resources={[ResourceType.Secrets]}>
        <SecretList />
      </GetResources>
    </>
  )
}

export default SecretsIndex
