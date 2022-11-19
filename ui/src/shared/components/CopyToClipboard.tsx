import React, {FunctionComponent, ChangeEvent} from 'react'

export const copy = async (text: string): Promise<boolean> => {
  let result: boolean

  try {
    if (navigator.clipboard) {
      // All browsers except IE
      await navigator.clipboard.writeText(text)
    } else {
      // IE
      const successful = document.execCommand('copy')
      if (!successful) {
        throw new Error('Copy command was unsuccessful using execCommand')
      }
    }

    result = true
  } catch (err) {
    result = false
  }

  return result
}

interface Props {
  text: string
  children: JSX.Element
}

const CopyToClipboard: FunctionComponent<Props> = ({
  text,
  children,
  ...props
}) => {
  const elmt = React.Children.only(children)

  const onClick = async (ev: ChangeEvent<HTMLInputElement>) => {
    const result = await copy(text)

    // Bypass onClick if it was present
    if (elmt && elmt.props && typeof elmt.props.onClick === 'function') {
      elmt.props.onClick(ev)
    }
  }

  return React.cloneElement(elmt, {...props, onClick})
}

export default CopyToClipboard
