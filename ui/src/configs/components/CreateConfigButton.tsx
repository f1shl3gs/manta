import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback} from 'react'
import {useNavigate} from 'react-router-dom'
import {useOrg} from 'src/organizations/selectors'

const CreateConfigButton: FunctionComponent = () => {
  const navigate = useNavigate()
  const {id: orgID} = useOrg()
  const create = useCallback(() => {
    navigate(`/orgs/${orgID}/data/config/new`)
  }, [navigate, orgID])

  return (
    <Button
      testID={'button-create-config'}
      text="Create Config"
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      onClick={create}
    />
  )
}

export default CreateConfigButton
