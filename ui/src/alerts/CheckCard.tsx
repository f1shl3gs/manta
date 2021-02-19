// Libraries
import React from 'react'

// Components
import {
  AlignItems,
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  FlexDirection,
  IconFont,
  JustifyContent,
  ResourceCard,
  SlideToggle,
} from '@influxdata/clockface'

// Types
import {Check} from '../types/Check'
import {relativeTimestampFormatter} from '../utils/relativeTimestampFormatter'
import LastRunStatus from './LastRunStatus'
import moment from 'moment'

interface Props {
  check: Check
}

const CheckCard: React.FC<Props> = props => {
  const {
    check: {id, name, desc, updated, status, conditions},
  } = props

  const contextMenu = () => (
    <Button
      icon={IconFont.Trash}
      text={'Delete'}
      color={ComponentColor.Danger}
      size={ComponentSize.ExtraSmall}
      onClick={() => console.log(`delete ${id}`)}
    />
  )

  return (
    <ResourceCard
      key={`check-id--${id}`}
      disabled={status === 'inactive'}
      direction={FlexDirection.Row}
      alignItems={AlignItems.Center}
      margin={ComponentSize.Large}
      contextMenu={contextMenu()}
    >
      <FlexBox
        direction={FlexDirection.Column}
        justifyContent={JustifyContent.Center}
        margin={ComponentSize.Medium}
        alignItems={AlignItems.FlexStart}
      >
        <SlideToggle
          active={status !== 'inactive'}
          size={ComponentSize.ExtraSmall}
          onChange={() => console.log('toggle')}
          style={{flexBasis: '16px'}}
        />
        <LastRunStatus lastRunStatus={'success'} />
      </FlexBox>

      <FlexBox
        direction={FlexDirection.Column}
        margin={ComponentSize.Small}
        alignItems={AlignItems.FlexStart}
      >
        <ResourceCard.EditableName
          name={name}
          noNameString={'Name this Check'}
          onClick={() => console.log('onClick')}
          onUpdate={v => console.log('v')}
        />
        <ResourceCard.EditableDescription
          description={desc || ''}
          placeholder={`Describe ${name}`}
          onUpdate={v => console.log('onUpdate')}
        />
        <ResourceCard.Meta>
          <>Last completed at {moment().format()}</>
          <>{relativeTimestampFormatter(updated, 'Last updated ')}</>
        </ResourceCard.Meta>
      </FlexBox>
    </ResourceCard>
  )
}

export default CheckCard
