import React, {FunctionComponent} from 'react'
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useOrganization} from '../organizations/useOrganizations'
import {useNotification} from '../shared/components/notifications/useNotification'
import {NotificationStyle} from '../types/Notification'
import {useNavigate} from 'react-router-dom'
import useFetch from '../shared/useFetch'

const CreateDashboardButton: FunctionComponent = () => {
  const {id: orgId} = useOrganization()
  const {notify} = useNotification()
  const navigate = useNavigate()
  const {run: create} = useFetch(`/api/v1/dashboards`, {
    method: 'POST',
    body: {
      orgId,
    },
    onError: err => {
      notify({
        icon: IconFont.AlertTriangle,
        style: NotificationStyle.Error,
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
      onClick={create}
    />
  )
}

export default CreateDashboardButton
