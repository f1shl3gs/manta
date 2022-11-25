// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ComponentColor,
  ComponentSize,
  Dropdown,
  IconFont,
} from '@influxdata/clockface'

interface Props {
  resourceType: string

  onSelectNew: () => void
  onSelectImport?: () => void
}

const AddResourceDropdown: FunctionComponent<Props> = props => {
  const {resourceType, onSelectNew, onSelectImport} = props

  const options = () => {
    const newOption = (
      <Dropdown.Item
        id="new"
        key="new"
        onClick={onSelectNew}
        value={`New ${resourceType}`}
        testID={'add-resource-dropdown--new'}
      >
        {`New ${resourceType}`}
      </Dropdown.Item>
    )

    if (!onSelectImport) {
      return [newOption]
    }

    return [
      newOption,
      <Dropdown.Item
        id="import"
        key="import"
        onClick={onSelectImport}
        value={`Import ${resourceType}`}
        testID={'add-resource-dropdown--import'}
      >
        {`Import ${resourceType}`}
      </Dropdown.Item>,
    ]
  }

  return (
    <Dropdown
      style={{width: 'fit-content'}}
      testID={'add-resource-dropdown'}
      button={(active, onClick) => (
        <Dropdown.Button
          testID={'add-resource-dropdown--button'}
          active={active}
          onClick={onClick}
          color={ComponentColor.Primary}
          size={ComponentSize.Small}
          icon={IconFont.Plus_New}
          style={{textTransform: 'uppercase', letterSpacing: '0.07em'}}
        >
          {`Create ${resourceType}`}
        </Dropdown.Button>
      )}
      menu={onCollapse => (
        <Dropdown.Menu
          onCollapse={onCollapse}
          testID={'add-resource-dropdown--menu'}
        >
          {options()}
        </Dropdown.Menu>
      )}
    >
      {options()}
    </Dropdown>
  )
}

export default AddResourceDropdown
