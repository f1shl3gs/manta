// Libraries
import React, {FunctionComponent, MouseEvent} from 'react'
import classnames from 'classnames'

// Types
import {ComponentColor, ComponentSize} from '@influxdata/clockface'

interface Props {
  text: string
  onClick: (ev: MouseEvent) => void
  color?: ComponentColor
  size?: ComponentSize
  testID?: string
}

const DashedButton: FunctionComponent<Props> = ({
  text,
  onClick,
  color = ComponentColor.Primary,
  size = ComponentSize.Medium,
  testID,
}) => {
  const classname = classnames('dashed-button', {
    [`dashed-button__${color}`]: color,
    [`dashed-button__${size}`]: size,
  })

  return (
    <button
      className={classname}
      onClick={onClick}
      type={'button'}
      data-testid={testID}
    >
      {text}
    </button>
  )
}

export default DashedButton
