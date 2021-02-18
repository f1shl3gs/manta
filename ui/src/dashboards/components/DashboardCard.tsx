import React from 'react'
import {ComponentColor, IconFont, ResourceCard} from '@influxdata/clockface'
import {useHistory} from 'react-router-dom'
import {useOrgID} from 'shared/useOrg'
import Context from 'components/context_menu/Context'

interface Props {
  id: string
  name: string
  desc: string
  updatedAt: string
  onDeleteDashboard: (id: string) => void
}

const DashboardCard: React.FC<Props> = props => {
  const {id, name, desc, updatedAt, onDeleteDashboard} = props
  const history = useHistory()
  const orgID = useOrgID()

  const contextMenu = (): JSX.Element => {
    return (
      <Context>
        <Context.Menu icon={IconFont.CogThick}>
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
        onUpdate={v => console.log('update dashboard name', v)}
        onClick={() => history.push(`/orgs/${orgID}/dashboards/${id}`)}
        name={name}
      />

      <ResourceCard.EditableDescription
        onUpdate={desc => console.log('update desc', desc)}
        description={desc}
        placeholder={`Describe ${name}`}
      />
      <ResourceCard.Meta>{`Last Modified: ${updatedAt}`}</ResourceCard.Meta>
    </ResourceCard>
  )
}

export default DashboardCard
