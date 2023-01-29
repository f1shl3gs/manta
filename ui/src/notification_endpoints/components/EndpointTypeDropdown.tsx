// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Dropdown} from '@influxdata/clockface'

// Types
import {NotificationEndpointType} from 'src/types/notificationEndpoints'

const types = [
  {
    name: 'HTTP',
    type: 'http',
  },
]

interface Props {
  selected: string
  onSelect: (type: NotificationEndpointType) => void
}

const EndpointTypeDropdown: FunctionComponent<Props> = ({
  selected,
  onSelect,
}) => {
  const items = types.map(({name, type}) => (
    <Dropdown.Item
      key={type}
      id={type}
      value={type}
      testID={`endpoint-dropdown--${type}`}
      onClick={onSelect}
    >
      {name}
    </Dropdown.Item>
  ))

  const selectedItem = types.find(item => item.type === selected)
  const button = (active, onClick) => (
    <Dropdown.Button active={active} onClick={onClick}>
      {selectedItem.name}
    </Dropdown.Button>
  )

  const menu = onCollapse => (
    <Dropdown.Menu onCollapse={onCollapse}>{items}</Dropdown.Menu>
  )

  return <Dropdown button={button} menu={menu} />
}

export default EndpointTypeDropdown
