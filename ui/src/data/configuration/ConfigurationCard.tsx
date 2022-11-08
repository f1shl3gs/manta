import React, {FunctionComponent} from 'react'
import {Configuration} from 'src/types/Configuration'
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import {fromNow} from 'src/shared/duration'
import Context from 'src/shared/components/context_menu/Context'
import {useNavigate} from 'react-router-dom'
import {useOrganization} from 'src/organizations/useOrganizations'

interface Props {
  configuration: Configuration
}

const ConfigurationCard: FunctionComponent<Props> = ({configuration}) => {
  const navigate = useNavigate()
  const {id: orgId} = useOrganization()

  const contextMenu = (): JSX.Element => (
    <Context>
      <Context.Menu
        icon={IconFont.Trash_New}
        color={ComponentColor.Danger}
        shape={ButtonShape.Square}
        testID={'configuration-card-context--delete'}
      >
        <Context.Item
          label={'Delete'}
          action={() => console.log('delete')}
          testID={'context_menu-delete'}
        />
      </Context.Menu>
    </Context>
  )

  return (
    <ResourceCard
      key={configuration.id}
      testID={'configration-card'}
      contextMenu={contextMenu()}
    >
      <ResourceCard.EditableName
        name={configuration.name}
        onUpdate={t => console.log(t)}
        onClick={() => {
          navigate(`/orgs/${orgId}/data/config/${configuration.id}`)
        }}
      />

      <ResourceCard.EditableDescription
        description={configuration.desc}
        placeholder={`Describe ${configuration.name}`}
        onUpdate={desc => console.log('update', desc)}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(configuration.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default ConfigurationCard
