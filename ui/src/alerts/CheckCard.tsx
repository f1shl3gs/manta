import React from 'react'
import {Check} from '../types/Check'
import {
  AlignItems,
  ComponentSize,
  FlexBox,
  FlexDirection,
  JustifyContent,
  ResourceCard,
  SlideToggle,
} from '@influxdata/clockface'

interface Props {
  check: Check
}

const CheckCard: React.FC<Props> = (props) => {
  const {
    check: {id, name, desc, updated, status, conditions},
  } = props

  /* todo: context Menu */
  return (
    <ResourceCard
      key={`check-id--${id}`}
      direction={FlexDirection.Row}
      alignItems={AlignItems.Center}
      margin={ComponentSize.Large}
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
      </FlexBox>
    </ResourceCard>
  )
}

export default CheckCard
