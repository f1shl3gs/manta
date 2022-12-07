// Libraries
import React, {FunctionComponent, useCallback} from 'react'

// Components
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import Context from 'src/shared/components/context_menu/Context'

// Hooks
import {useNavigate} from 'react-router-dom'
import {fromNow} from 'src/shared/utils/duration'
import {useOrg} from 'src/organizations/selectors'
import {useDispatch, useSelector} from 'react-redux'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'
import {Dashboard} from 'src/types/dashboards'

// Constants
import {
  PARAMS_INTERVAL,
  PARAMS_SHOW_VARIABLES_CONTROLS,
  PARAMS_TIME_RANGE_LOW,
  PARAMS_TIME_RANGE_TYPE,
} from 'src/shared/constants/timeRange'

// Actions
import {
  cloneDashboard,
  deleteDashboard,
  updateDashboard,
} from 'src/dashboards/actions/thunks'

// Selectors
import {getByID} from 'src/resources/selectors'

interface Props {
  id: string
}

const DashboardCard: FunctionComponent<Props> = props => {
  const {id: orgID} = useOrg()
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const dashboard = useSelector((state: AppState) => {
    return getByID<Dashboard>(state, ResourceType.Dashboards, props.id)
  })

  const handleDelete = () => {
    dispatch(deleteDashboard(dashboard.id, dashboard.name))
  }

  const handleUpdate = upd => {
    dispatch(updateDashboard(dashboard.id, upd))
  }

  const handleClone = useCallback(() => {
    const {id, name} = dashboard
    dispatch(cloneDashboard(id, name))
  }, [dashboard, dispatch])

  const handleClick = () => {
    navigate(
      `/orgs/${orgID}/dashboards/${dashboard.id}?${new URLSearchParams({
        [PARAMS_INTERVAL]: '15s',
        [PARAMS_TIME_RANGE_LOW]: 'now() - 1h',
        [PARAMS_TIME_RANGE_TYPE]: 'selectable-duration',
        [PARAMS_SHOW_VARIABLES_CONTROLS]: 'true',
      })}`
    )
  }

  const handleExport = useCallback(() => {
    navigate(`${window.location.pathname}/${dashboard.id}/export`)
  }, [dashboard, navigate])

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
          action={handleDelete}
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
        onUpdate={name => handleUpdate({name})}
        onClick={handleClick}
      />

      <ResourceCard.EditableDescription
        testID={'dashboard-editable-desc'}
        description={dashboard.desc}
        placeholder={`Describe ${dashboard.name}`}
        onUpdate={desc => handleUpdate({desc})}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(dashboard.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default DashboardCard
