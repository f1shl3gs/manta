// Libraries
import {Dropdown} from '@influxdata/clockface'
import React, {FunctionComponent} from 'react'

const methods = ['POST']

interface Props {
  selected: string
  onSelect: (method: string) => void
}

const MethodDropdown: FunctionComponent<Props> = ({selected, onSelect}) => {
  const items = methods.map(method => (
    <Dropdown.Item
      key={method}
      id={method}
      value={method}
      testID={`http-method-dropdown--${method}`}
      onClick={onSelect}
    >
      {method}
    </Dropdown.Item>
  ))

  const selectedMethod = methods.find(m => m == selected)
  const button = (active, onClick) => (
    <Dropdown.Button
      testID={`http-method--dropdown--button`}
      active={active}
      onClick={onClick}
    >
      {selectedMethod}
    </Dropdown.Button>
  )

  const menu = onCollapse => (
    <Dropdown.Menu onCollapse={onCollapse}>{items}</Dropdown.Menu>
  )

  return (
    <Dropdown button={button} menu={menu} testID={'http-method--dropdown'} />
  )
}

export default MethodDropdown
