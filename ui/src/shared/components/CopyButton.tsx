// Libraries
import React from 'react'

// Components
import {
  Button,
  ButtonShape,
  ComponentColor,
  ComponentSize,
  IconFont,
  Notification,
} from '@influxdata/clockface'
//@ts-ignore
import CopyToClipboard from 'react-copy-to-clipboard'

interface Props {
  shape: ButtonShape
  icon?: IconFont
  buttonText: string
  textToCopy: string
  contentName: string
  size: ComponentSize
  color: ComponentColor
  onCopyText?: (text: string, status: boolean) => Notification
  onClick?: () => void
}

const CopyButton: React.FC<Props> = props => {
  const {
    textToCopy,
    color,
    size,
    icon,
    shape,
    buttonText,
    onCopyText,
    onClick,
  } = props

  const handleCopy = (copiedText: string, isSuccessful: boolean) => {
    if (onClick) {
      onClick()
    }

    if (onCopyText) {
      onCopyText(copiedText, isSuccessful)
    }

    console.log('todo notification')
  }

  return (
    <CopyToClipboard text={textToCopy} onCopy={handleCopy}>
      <Button
        shape={shape}
        icon={icon}
        size={size}
        color={color}
        titleText={buttonText}
        text={buttonText}
        // onClick={undefined}
      />
    </CopyToClipboard>
  )
}

export default CopyButton
