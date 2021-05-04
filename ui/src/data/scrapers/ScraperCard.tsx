// Libraries
import React from 'react'
import {useHistory} from 'react-router-dom'
import moment from 'moment'

// Components
import {
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  IconFont,
  InfluxColors,
  Label,
  ResourceCard,
} from '@influxdata/clockface'

// Hooks
import {useScrapers} from './useScrapers'
import {useOrgID} from '../../shared/useOrg'

// Types
import {Scraper} from '../../types/scrapers'

interface Props {
  scraper: Scraper
}

const ScraperCard: React.FC<Props> = props => {
  const {
    scraper: {id, name, desc, updated, labels},
  } = props
  const {onRemove, onNameUpdate, onDescUpdate} = useScrapers()
  const history = useHistory()
  const orgID = useOrgID()

  const context = (): JSX.Element => (
    <FlexBox margin={ComponentSize.Small}>
      <FlexBox.Child>
        <Button
          icon={IconFont.Trash}
          text={'Delete'}
          size={ComponentSize.ExtraSmall}
          color={ComponentColor.Danger}
          onClick={() => onRemove(id)}
        />
      </FlexBox.Child>
    </FlexBox>
  )

  return (
    <ResourceCard key={id} contextMenu={context()}>
      <ResourceCard.EditableName
        name={name}
        onUpdate={name => onNameUpdate(id, name)}
        onClick={() => {
          history.push(`/orgs/${orgID}/data/scrapers/${id}`)
        }}
      />

      <ResourceCard.EditableDescription
        description={desc}
        onUpdate={desc => onDescUpdate(id, desc)}
      />

      <ResourceCard.Meta>
        <span>updated: {moment(updated).fromNow()}</span>
      </ResourceCard.Meta>

      <div className={'inline-labels--container'}>
        {Object.keys(labels).map(key => (
          <Label
            id={key}
            key={key}
            name={key}
            color={InfluxColors.Abyss}
            description={''}
          />
        ))}
      </div>
    </ResourceCard>
  )
}

export default ScraperCard
