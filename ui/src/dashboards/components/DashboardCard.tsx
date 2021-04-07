// Libraries
import React, {useCallback} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import Context from 'components/context_menu/Context'

// Hooks
import {useOrgID} from 'shared/useOrg'
import {useDashboards} from '../useDashboards'

interface Props {
  id: string
  name: string
  desc: string
  updatedAt: string
  // onNameUpdate: (name: string) => void
  onDeleteDashboard: (id: string) => void
}

const DashboardCard: React.FC<Props> = props => {
  const {id, name, desc, updatedAt, onDeleteDashboard} = props
  const history = useHistory()
  const orgID = useOrgID()
  const {updateDashboard} = useDashboards()

  const onNameUpdate = useCallback(
    (name: string) => {
      updateDashboard(id, {
        name,
      })
    },
    [id, updateDashboard]
  )

  const onDescUpdate = useCallback(
    (desc: string) => {
      updateDashboard(id, {
        desc,
      })
    },
    [id, updateDashboard]
  )

  const contextMenu = (): JSX.Element => {
    return (
      <Context>
        <Context.Menu
          icon={IconFont.CogThick}
          color={ComponentColor.Default}
          shape={ButtonShape.Square}
          testID={'dashboard-card-context--export'}
        >
          <Context.Item
            label={'Export'}
            action={value => console.log('export action', value)}
            testID={'context_menu-export'}
          />
        </Context.Menu>
        <Context.Menu icon={IconFont.Trash} color={ComponentColor.Danger}>
          <Context.Item
            label={'Delete'}
            action={value => onDeleteDashboard(id)}
          />
        </Context.Menu>
      </Context>
    )
  }

  return (
    <ResourceCard key={`dashboard-id--${id}`} contextMenu={contextMenu()}>
      <ResourceCard.EditableName
        onUpdate={onNameUpdate}
        onClick={() =>
          history.push(
            `/orgs/${orgID}/dashboards/${id}?${new URLSearchParams({
              _interval: '15s',
              _lower: 'now() - 1h',
              _type: 'selectable-duration',
            }).toString()}`
          )
        }
        name={name}
      />

      <ResourceCard.EditableDescription
        onUpdate={onDescUpdate}
        description={desc}
        placeholder={`Describe ${name}`}
      />
      <ResourceCard.Meta>{`Last Modified: ${updatedAt}`}</ResourceCard.Meta>
    </ResourceCard>
  )
}

export default DashboardCard
