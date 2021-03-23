// Libraries
import React from 'react'
import classnames from 'classnames'

// Components
import {ComponentColor, ComponentSize} from '@influxdata/clockface'

interface Props {
  text: string
  onClick: (e: MouseEvent) => void
  color?: ComponentColor
  size?: ComponentSize
  testID?: string
}

const DashedButton: React.FC<Props> = props => {
  const {
    text,
    onClick,
    color = ComponentColor.Primary,
    size = ComponentSize.Medium,
    testID = 'dashed-button',
  } = props

  const classname = classnames('dashed-button', {
    [`dashed-button__${color}`]: color,
    [`dashed-button__${size}`]: size,
  })

  return (
    <button
      className={classname}
      // @ts-ignore
      onClick={onClick}
      type={'button'}
      data-testid={testID}
    >
      {text}
    </button>
  )
}

export default DashedButton
