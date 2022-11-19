// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

// Hooks
import {
  useNotification,
  defaultErrorNotification,
} from 'src/shared/components/notifications/useNotification'
import {useNavigate} from 'react-router-dom'
import {useOrganization} from 'src/organizations/useOrganizations'
import useFetch from 'src/shared/useFetch'

const CreateDashboardButton: FunctionComponent = () => {
  const {id: orgId} = useOrganization()
  const {notify} = useNotification()
  const navigate = useNavigate()
  const {run: create} = useFetch(`/api/v1/dashboards`, {
    method: 'POST',
    body: {
      orgId,
      cells: [],
    },
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Create new dashboard failed\n${err}`,
      })
    },
    onSuccess: dashboard => {
      navigate(`${window.location.pathname}/${dashboard.id}`)
    },
  })

  return (
    <Button
      testID={'button-create-dashboard'}
      text="Create Dashboard"
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      onClick={() => create()}
    />
  )
}

export default CreateDashboardButton
