import React from 'react'

// Components
import {ComponentSize, SelectableCard, SquareGrid} from '@influxdata/clockface'

// Graphics
import placeholderLogo from 'plugins/graphics/placeholderLogo.svg'

// Hooks
import {useHistory} from 'react-router-dom'
import {useOrgID} from 'shared/useOrg'

// Constants
import {PluginItem} from './constants/Plugins'

type Props = PluginItem

const PluginCard: React.FC<Props> = (props) => {
  const {id, name, url, image} = props
  const history = useHistory()
  const orgID = useOrgID()

  const handleClick = () => {
    history.push(`/orgs/${orgID}/${url}`)
  }

  return (
    <SquareGrid.Card key={id}>
      <SelectableCard
        id={id}
        formName={'plugin-cards'}
        label={name}
        selected={false}
        onClick={handleClick}
        fontSize={ComponentSize.ExtraSmall}
        className={'write-data--item'}
      >
        <div className={'write-data--item-thumb'}>
          <img src={image ? image : placeholderLogo} />
        </div>
      </SelectableCard>
    </SquareGrid.Card>
  )
}

export default PluginCard
