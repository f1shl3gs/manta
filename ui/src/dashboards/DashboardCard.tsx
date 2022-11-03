import React, {FunctionComponent} from 'react'
import {Dashboard} from 'types/Dashboard'
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import Context from 'shared/components/context_menu/Context'
import {useNavigate} from 'react-router-dom'
import {fromNow} from 'shared/duration'
import {useOrganization} from '../organizations/useOrganizations'
import {
  PARAMS_INTERVAL,
  PARAMS_SHOW_VARIABLES_CONTROLS,
  PARAMS_TIME_RANGE_LOW,
  PARAMS_TIME_RANGE_TYPE,
} from './constants'
import useFetch from 'shared/useFetch'
import {useNotification} from '../shared/components/notifications/useNotification';
import {NotificationStyle} from '../types/Notification';

interface Props {
  dashboard: Dashboard
  onDelete: () => void
}

const DashboardCard: FunctionComponent<Props> = props => {
  const {dashboard, onDelete} = props
  const navigate = useNavigate()
  const {id: orgId} = useOrganization()
  const {notify} = useNotification()
  const {run: deleteDashboard} = useFetch(
    `/api/v1/dashboards/${dashboard.id}`,
    {
      method: 'DELETE',
      onError: err => {
        notify({
          style: NotificationStyle.Error,
          icon: IconFont.AlertTriangle,
          message: `Delete dashboard ${dashboard.name} failed, err: ${err}`
        })
      },
      onSuccess: _ => {
        console.log('success')

        onDelete()
        notify({
          style: NotificationStyle.Success,
          icon: IconFont.CheckMark_New,
          message: `Delete dashboard ${dashboard.name} success`
        })
      },
    }
  )

  const contextMenu = (): JSX.Element => (
    <Context>
      <Context.Menu
        icon={IconFont.CogOutline_New}
        color={ComponentColor.Default}
        shape={ButtonShape.Square}
        testID="dashboard-card-context--export"
      >
        <Context.Item
          label="Export"
          action={value => console.log('export', value)}
          testID="context_menu-export"
        />
      </Context.Menu>

      <Context.Menu
        icon={IconFont.Trash_New}
        color={ComponentColor.Danger}
        shape={ButtonShape.Square}
        testID="dashboard-card-context--delete"
      >
        <Context.Item
          label="Delete"
          action={deleteDashboard}
          testID="context_menu-delete"
        />
      </Context.Menu>
    </Context>
  )

  return (
    <ResourceCard key={dashboard.id} contextMenu={contextMenu()} testID={'dashboard-card'}>
      <ResourceCard.EditableName
        testID={'dashboard-editable-name'}
        name={dashboard.name}
        onUpdate={v => console.log('onupdate', v)}
        onClick={() => {
          navigate(
            `/orgs/${orgId}/dashboards/${dashboard.id}?${new URLSearchParams({
              [PARAMS_INTERVAL]: '15s',
              [PARAMS_TIME_RANGE_LOW]: 'now() - 1h',
              [PARAMS_TIME_RANGE_TYPE]: 'selectable-duration',
              [PARAMS_SHOW_VARIABLES_CONTROLS]: 'true',
            })}`
          )
        }}
      />

      <ResourceCard.EditableDescription
        description={dashboard.desc}
        placeholder={`Describe ${dashboard.name}`}
        onUpdate={() => console.log('update desc')}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(dashboard.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default DashboardCard
