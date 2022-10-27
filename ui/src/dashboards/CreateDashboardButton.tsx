import React, {FunctionComponent, useCallback} from 'react'
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useOrganization} from '../organizations/useOrganizations'
import {useNotification} from '../shared/components/notifications/useNotification'
import {NotificationStyle} from '../types/Notification'
import {useNavigate} from 'react-router-dom'

const CreateDashboardButton: FunctionComponent = () => {
  const {id: orgId} = useOrganization()
  const {notify} = useNotification()
  const navigate = useNavigate()

  const create = useCallback(() => {
    fetch(`/api/v1/dashboards`, {
      method: 'POST',
      body: JSON.stringify({
        orgId,
      }),
    })
      .then(resp => {
        if (resp.status !== 201) {
          notify({
            icon: IconFont.AlertTriangle,
            style: NotificationStyle.Error,
            message: 'Create new dashboard failed',
          })

          return
        }

        return resp.json()
      })
      .then(dashboard => {
        navigate(`${window.location.pathname}/${dashboard.id}`)
      })
      .catch(err =>
        notify({
          icon: IconFont.AlertTriangle,
          style: NotificationStyle.Error,
          message: `Create new dashboard failed\n${err}`,
        })
      )
  }, [])

  return (
    <Button
      text="Create Dashboard"
      icon={IconFont.Plus_New}
      color={ComponentColor.Primary}
      onClick={create}
    />
  )
}

export default CreateDashboardButton
