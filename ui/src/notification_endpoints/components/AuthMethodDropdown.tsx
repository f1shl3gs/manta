// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Dropdown} from '@influxdata/clockface'

// Types
import {HTTPAuthMethod} from 'src/types/notificationEndpoints'

const types = ['none', 'basic', 'bearer']

interface Props {
  selected: HTTPAuthMethod
  onSelect: (authMethod: HTTPAuthMethod) => void
}

const AuthMethodDropdown: FunctionComponent<Props> = ({selected, onSelect}) => {
  const items = types.map(type => (
    <Dropdown.Item
      key={type}
      id={type}
      value={type}
      testID={`http-auth-method-${type}`}
      onClick={onSelect}
    >
      {type}
    </Dropdown.Item>
  ))

  const button = (active, onClick) => (
    <Dropdown.Button
      active={active}
      onClick={onClick}
      testID={'http-auth-method--dropdown-button'}
    >
      {selected}
    </Dropdown.Button>
  )

  const menu = onCollapse => (
    <Dropdown.Menu onCollapse={onCollapse}>{items}</Dropdown.Menu>
  )

  return (
    <Dropdown
      button={button}
      menu={menu}
      testID={'http-auth-method--dropdown'}
    />
  )
}

export default AuthMethodDropdown
