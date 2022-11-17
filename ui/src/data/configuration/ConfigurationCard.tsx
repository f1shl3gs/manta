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
import useFetch from 'src/shared/useFetch'
import {useResources} from 'src/shared/components/GetResources'
import {
  defaultSuccessNotification,
  defaultErrorNotification,
  useNotification,
} from 'src/shared/components/notifications/useNotification'

interface Props {
  configuration: Configuration
}

const ConfigurationCard: FunctionComponent<Props> = ({configuration}) => {
  const navigate = useNavigate()
  const {id: orgId} = useOrganization()
  const {reload} = useResources()
  const {notify} = useNotification()
  const {run: deleteConfig} = useFetch(
    `/api/v1/configurations/${configuration.id}`,
    {
      method: 'DELETE',
      onSuccess: _ => {
        notify({
          ...defaultSuccessNotification,
          message: `Delete configuration ${configuration.name} success`,
        })

        reload()
      },
      onError: err => {
        notify({
          ...defaultErrorNotification,
          message: `Delete configuration ${configuration.name} failed, ${err}`,
        })
      },
    }
  )
  const {run: update} = useFetch(`/api/v1/configurations/${configuration.id}`, {
    method: 'PATCH',
    onSuccess: _ => {
      reload()
    },
  })

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
          action={() => deleteConfig()}
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
        onUpdate={name => update({name})}
        onClick={() => {
          navigate(`/orgs/${orgId}/data/config/${configuration.id}`)
        }}
      />

      <ResourceCard.EditableDescription
        description={configuration.desc}
        placeholder={`Describe ${configuration.name}`}
        onUpdate={desc => update({desc})}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(configuration.updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default ConfigurationCard
