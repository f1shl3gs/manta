// Libraries
import React, {FunctionComponent, RefObject, useCallback, useRef} from 'react'

// Components
import {
  Appearance,
  ButtonShape,
  ComponentColor,
  ComponentSize,
  ConfirmationButton,
  FlexBox,
  IconFont,
  List,
  Popover,
  ResourceCard,
  SquareButton,
} from '@influxdata/clockface'

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

const fontWeight = {fontWeight: '500px'}
const minWidth = {minWidth: '165px'}

interface Props {
  id: string
}

const DashboardCard: FunctionComponent<Props> = props => {
  const {id: orgID} = useOrg()
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const settingRef = useRef<RefObject<HTMLButtonElement>>(null)
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

  const contextMenu = (ref): JSX.Element => (
    <FlexBox margin={ComponentSize.ExtraSmall}>
      <ConfirmationButton
        color={ComponentColor.Colorless}
        icon={IconFont.Trash_New}
        shape={ButtonShape.Square}
        size={ComponentSize.ExtraSmall}
        confirmationLabel={'Delete this dashboard'}
        confirmationButtonText={'Confirm'}
        onConfirm={handleDelete}
        testID={'dashboard-card-context--delete'}
      />

      <SquareButton
        ref={ref}
        size={ComponentSize.ExtraSmall}
        icon={IconFont.CogSolid_New}
        color={ComponentColor.Colorless}
        testID={'dashboard-card-context-menu'}
      />

      <Popover
        appearance={Appearance.Outline}
        enableDefaultStyles={false}
        style={minWidth}
        contents={_ => (
          <List>
            <List.Item
              onClick={handleClone}
              size={ComponentSize.ExtraSmall}
              style={fontWeight}
              testID={'dashboard-card-context--clone'}
            >
              Clone
            </List.Item>

            <List.Item
              onClick={handleExport}
              size={ComponentSize.ExtraSmall}
              style={fontWeight}
              testID={'dashboard-card-context--export'}
              >
              Export
            </List.Item>
          </List>
        )}
        triggerRef={ref}
      />
    </FlexBox>
  )

  return (
    <ResourceCard
      key={dashboard.id}
      contextMenu={contextMenu(settingRef)}
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
