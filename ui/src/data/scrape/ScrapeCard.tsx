import React, {FunctionComponent} from 'react'
import {Scrape} from 'src/types/Scrape'
import {
  ButtonShape,
  ComponentColor,
  IconFont,
  ResourceCard,
} from '@influxdata/clockface'
import useFetch from 'src/shared/useFetch'
import {useNavigate} from 'react-router-dom'
import {fromNow} from 'src/shared/duration'
import Context from 'src/shared/components/context_menu/Context'
import {useResources} from 'src/shared/components/GetResources'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'src/shared/components/notifications/useNotification'

interface Props {
  scrape: Scrape
}

const ScrapeCard: FunctionComponent<Props> = props => {
  const {id, name, desc, updated} = props.scrape
  const navigate = useNavigate()
  const {reload} = useResources()
  const {notify} = useNotification()
  const {run: patchScrape} = useFetch(`/api/v1/scrapes/${id}`, {method: 'POST'})
  const {run: deleteScrape} = useFetch(`/api/v1/scrapes/${id}`, {
    method: 'DELETE',
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Delete Scrape ${name} failed, ${err}`,
      })
    },
    onSuccess: _ => {
      reload()

      notify({
        ...defaultSuccessNotification,
        message: `Delete Scrape ${name} successful`,
      })
    },
  })

  const handleNameUpdate = (name: string): void => {
    patchScrape({name})
  }
  const handleDescUpdate = (desc: string): void => {
    patchScrape({desc})
  }

  const handleNameClick = (): void => {
    navigate(`${window.location.pathname}/${id}`)
  }

  const handleDelete = (): void => {
    deleteScrape()
  }

  const context_menu = (): JSX.Element => (
    <Context>
      <Context.Menu
        icon={IconFont.Trash_New}
        color={ComponentColor.Danger}
        shape={ButtonShape.Square}
        testID={'scrape-card-context--delete'}
      >
        <Context.Item label={'Delete'} action={handleDelete} />
      </Context.Menu>
    </Context>
  )

  return (
    <ResourceCard key={id} testID={'scrape-card'} contextMenu={context_menu()}>
      <ResourceCard.EditableName
        onUpdate={handleNameUpdate}
        name={name}
        onClick={handleNameClick}
      />

      <ResourceCard.EditableDescription
        description={desc}
        onUpdate={handleDescUpdate}
        placeholder={`Describe this ${name}`}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(updated)}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default ScrapeCard
