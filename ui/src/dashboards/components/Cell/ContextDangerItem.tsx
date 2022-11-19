// Libraries
import React, {FunctionComponent, useState} from 'react'

// Components
import {
  Button,
  ButtonShape,
  ComponentColor,
  ComponentSize,
  Icon,
  IconFont,
} from '@influxdata/clockface'

interface Props {
  label: string
  icon?: IconFont
  testID?: string
  onClick: () => void
  onHide?: () => void
  confirmationText?: string
}

const ContextDangerItem: FunctionComponent<Props> = ({
  icon = IconFont.Trash_New,
  label,
  testID,
  onHide,
  onClick,
  confirmationText,
}) => {
  const [confirming, setConfirmationState] = useState(false)

  const toggleConfirmationState = (): void => {
    setConfirmationState(true)
  }

  const handleClick = (): void => {
    if (onHide) {
      onHide()
    }

    onClick()
  }

  if (confirming) {
    return (
      <div className={'cell--context-item cell--context-item__confirm'}>
        <Button
          color={ComponentColor.Danger}
          text={confirmationText}
          onClick={handleClick}
          size={ComponentSize.ExtraSmall}
          shape={ButtonShape.StretchToFit}
          testID={`${testID}-confirm`}
        />
      </div>
    )
  }

  return (
    <div
      className={'cell--context-item cell--context-item__danger'}
      onClick={toggleConfirmationState}
      data-testid={testID}
    >
      <Icon glyph={icon} />
      {label}
    </div>
  )
}

export default ContextDangerItem
