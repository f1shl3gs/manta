import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback} from 'react'
import {useNavigate} from 'react-router-dom'
import {useOrganization} from 'src/organizations/useOrganizations'

const CreateConfigurationButton: FunctionComponent = () => {
  const navigate = useNavigate()
  const {id: orgId} = useOrganization()
  const create = useCallback(() => {
    navigate(`/orgs/${orgId}/data/config/new`)
  }, [navigate, orgId])

  return (
    <Button
      testID={'button-create-configuration'}
      text="Create Configuration"
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      onClick={create}
    />
  )
}

export default CreateConfigurationButton
