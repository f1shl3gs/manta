import React, {FunctionComponent, useCallback} from 'react'
import {Dashboard} from 'src/types/Dashboard'
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import Context from 'src/shared/components/context_menu/Context'
import {useNavigate} from 'react-router-dom'
import {fromNow} from 'src/utils/duration'
import {useOrganization} from 'src/organizations/useOrganizations'
import useFetch from 'src/shared/useFetch'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'
import {
  PARAMS_INTERVAL,
  PARAMS_TIME_RANGE_LOW,
  PARAMS_TIME_RANGE_TYPE,
  PARAMS_SHOW_VARIABLES_CONTROLS,
} from 'src/constants/timeRange'

interface Props {
  dashboard: Dashboard
  reload: () => void
}

const DashboardCard: FunctionComponent<Props> = props => {
  const {dashboard, reload} = props
  const navigate = useNavigate()
  const {id: orgId} = useOrganization()
  const notify = useNotify()
  const {run: deleteDashboard} = useFetch(
    `/api/v1/dashboards/${dashboard.id}`,
    {
      method: 'DELETE',
      onError: err => {
        notify({
          ...defaultErrorNotification,
          message: `Delete dashboard ${dashboard.name} failed, err: ${err}`,
        })
      },
      onSuccess: _ => {
        reload()

        notify({
          ...defaultSuccessNotification,
          message: `Delete dashboard ${dashboard.name} success`,
        })
      },
    }
  )
  const {run: update} = useFetch(`/api/v1/dashboards/${dashboard.id}`, {
    method: 'PATCH',
    onSuccess: reload,
  })
  const handleExport = useCallback(() => {
    navigate(`${window.location.pathname}/${dashboard.id}/export`)
  }, [dashboard, navigate])

  const {run: create} = useFetch(`/api/v1/dashboards`, {
    method: 'POST',
    onSuccess: _ => {
      reload()
    },
  })
  const handleClone = useCallback(() => {
    create({
      ...dashboard,
      name: `${dashboard.name} (Clone)`,
      orgID: orgId,
    })
  }, [create, dashboard, orgId])

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
          action={handleExport}
          testID="context_menu-export"
        />

        <Context.Item
          label="Clone"
          action={handleClone}
          testID={'context_menu-clone'}
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
    <ResourceCard
      key={dashboard.id}
      contextMenu={contextMenu()}
      testID={'dashboard-card'}
    >
      <ResourceCard.EditableName
        testID={'dashboard-editable-name'}
        buttonTestID={'dashboard-editable-name--button'}
        inputTestID={'dashboard-editable-name--input'}
        name={dashboard.name}
        onUpdate={name => update({name})}
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
        testID={'dashboard-editable-desc'}
        description={dashboard.desc}
        placeholder={`Describe ${dashboard.name}`}
        onUpdate={desc => update({desc})}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(dashboard.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default DashboardCard
