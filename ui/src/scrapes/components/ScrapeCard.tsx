// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  IconFont,
  InfluxColors,
  Label,
  ResourceCard,
} from '@influxdata/clockface'
import Context from 'src/shared/components/context_menu/Context'

// Types
import {Scrape} from 'src/types/scrape'

// Hooks
import {useNavigate} from 'react-router-dom'

// Utils
import {fromNow} from 'src/shared/utils/duration'

// Constants
import {deleteScrape, updateScrape} from 'src/scrapes/actions/thunks'

interface OwnProps {
  scrape: Scrape
}

const mdtp = {
  updateScrape,
  deleteScrape,
}

const connector = connect(null, mdtp)

type Props = OwnProps & ConnectedProps<typeof connector>

const ScrapeCard: FunctionComponent<Props> = ({
  scrape,
  deleteScrape,
  updateScrape,
}) => {
  const {id, name, desc, updated, labels} = scrape
  const navigate = useNavigate()

  const handleNameUpdate = (name: string): void => {
    updateScrape(id, {name})
  }
  const handleDescUpdate = (desc: string): void => {
    updateScrape(id, {desc})
  }

  const handleNameClick = (): void => {
    navigate(`${window.location.pathname}/${id}`)
  }

  const handleDelete = (): void => {
    deleteScrape(id)
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
        {Object.keys(labels).map(key => (
          <Label
            id={key}
            key={key}
            color={InfluxColors.Ocean}
            size={ComponentSize.ExtraSmall}
            name={`${key}=${labels[key]}`}
            description={labels[key]}
          />
        ))}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default connector(ScrapeCard)
