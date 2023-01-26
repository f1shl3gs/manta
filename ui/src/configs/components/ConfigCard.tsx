import React, {FunctionComponent} from 'react'
import {Config} from 'src/types/config'
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
} from 'src/configs/actions/thunks'

interface Props {
  config: Config
}

const ConfigCard: FunctionComponent<Props> = ({config}) => {
  const {id} = config
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
        testID={'config-card-context--delete'}
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
      key={config.id}
      testID={'config-card'}
      contextMenu={contextMenu()}
    >
      <ResourceCard.EditableName
        name={config.name}
        onUpdate={name => handleUpdate({name})}
        onClick={() => {
          navigate(`/orgs/${orgID}/data/config/${config.id}`)
        }}
      />

      <ResourceCard.EditableDescription
        description={config.desc}
        placeholder={`Describe ${config.name}`}
        onUpdate={desc => handleUpdate({desc})}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(config.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default ConfigCard
