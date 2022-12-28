// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  ConfirmationButton,
  FlexBox,
  IconFont,
  SlideToggle,
} from '@influxdata/clockface'

// Types
import {CheckStatus} from 'src/types/checks'

interface Props {
  activeStatus: CheckStatus
  onActiveToggle: () => void
  onDelete: () => void
}

const CheckCardContext: FunctionComponent<Props> = ({
  activeStatus,
  onActiveToggle,
  onDelete,
}) => {
  return (
    <FlexBox margin={ComponentSize.Medium}>
      <SlideToggle
        onChange={onActiveToggle}
        active={activeStatus === 'active'}
        style={{height: '16px'}}
      />

      <ConfirmationButton
        color={ComponentColor.Colorless}
        icon={IconFont.Trash_New}
        shape={ButtonShape.Square}
        size={ComponentSize.ExtraSmall}
        testID={'context-delete--button'}
        confirmationLabel={'Delete this check'}
        confirmationButtonText={'Confirm'}
        onConfirm={onDelete}
      />
    </FlexBox>
  )
}

export default CheckCardContext
