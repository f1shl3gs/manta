import React, {FunctionComponent} from 'react'
import {Icon, IconFont} from '@influxdata/clockface'

interface Props {
  label: string
  testID?: string
  icon: IconFont
  onClick: () => void
  onHide?: () => void
}

const ContextItem: FunctionComponent<Props> = ({
  label,
  testID,
  icon,
  onClick,
  onHide,
}) => {
  const handleClick = (): void => {
    if (onHide) {
      onHide()
    }

    onClick()
  }

  return (
    <div
      className={'cell--context-item'}
      onClick={handleClick}
      data-testid={testID}
    >
      <Icon glyph={icon} />
      {label}
    </div>
  )
}

export default ContextItem
