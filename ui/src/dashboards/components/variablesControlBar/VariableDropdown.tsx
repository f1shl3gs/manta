import React, {useState} from 'react'
import {
  ComponentStatus,
  Dropdown,
  DropdownMenuTheme,
} from '@influxdata/clockface'

interface Props {
  variableID: string
  testID?: string
  onSelect?: () => void
}

const VariableDropdown: React.FC<Props> = props => {
  const values = ['1', '2', '3']

  const [selectedValue, setSelectedValue] = useState(() => {
    return values[0]
  })
  const {testID = 'variable-dropdown'} = props
  const dropdownStatus =
    values.length === 0 ? ComponentStatus.Disabled : ComponentStatus.Default

  const longestItemWidth = Math.floor(
    values.reduce(function (a, b) {
      return a.length > b.length ? a : b
    }, '').length * 9.5
  )

  const widthLength = Math.max(140, longestItemWidth)

  return (
    <div className={'variable-dropdown'}>
      <Dropdown
        style={{width: `${140}px`}}
        className={'variable-dropdown--dropdown'}
        testID={testID}
        button={(active, onClick) => (
          <Dropdown.Button
            active={active}
            onClick={onClick}
            testID={'variable-dropdown--button'}
            status={dropdownStatus}
          >
            {/* todo: handle loading state */}
            {selectedValue}
          </Dropdown.Button>
        )}
        menu={onCollapse => (
          <Dropdown.Menu
            style={{width: `${widthLength}px`}}
            onCollapse={onCollapse}
            theme={DropdownMenuTheme.Amethyst}
          >
            {values.map(val => (
              <Dropdown.Item
                key={val}
                id={val}
                value={val}
                onClick={setSelectedValue}
                selected={val === selectedValue}
                testID={'variable-dropdown--item'}
                className={'variable-dropdown--item'}
              >
                {val}
              </Dropdown.Item>
            ))}
          </Dropdown.Menu>
        )}
      />
    </div>
  )
}

export default VariableDropdown
