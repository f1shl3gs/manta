import React, {FunctionComponent} from 'react'
import {Configuration} from 'src/types/configuration'
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import {fromNow} from 'src/shared/utils/duration'
import Context from 'src/shared/components/context_menu/Context'
import {useDispatch} from 'react-redux'
import {useNavigate, useParams} from 'react-router-dom'
import {
  ConfigUpdate,
  deleteConfig,
  updateConfig,
} from 'src/configurations/actions/thunk'

interface Props {
  configuration: Configuration
}

const ConfigurationCard: FunctionComponent<Props> = ({configuration}) => {
  const {id} = configuration
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const {orgID} = useParams()

  const handleDelete = () => {
    console.log('delete', id)

    dispatch(deleteConfig(id))
  }

  const handleUpdate = (upd: ConfigUpdate) => {
    dispatch(updateConfig(id, upd))
  }

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
          action={handleDelete}
          testID={'context_menu-delete'}
        />
      </Context.Menu>
    </Context>
  )

  return (
    <ResourceCard
      key={configuration.id}
      testID={'configuration-card'}
      contextMenu={contextMenu()}
    >
      <ResourceCard.EditableName
        name={configuration.name}
        onUpdate={name => handleUpdate({name})}
        onClick={() => {
          navigate(`/orgs/${orgID}/data/config/${configuration.id}`)
        }}
      />

      <ResourceCard.EditableDescription
        description={configuration.desc}
        placeholder={`Describe ${configuration.name}`}
        onUpdate={desc => handleUpdate({desc})}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(configuration.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default ConfigurationCard
