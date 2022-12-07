// Libraries
import React, {FunctionComponent} from 'react'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {
  ComponentStatus,
  Dropdown,
  DropdownMenuTheme,
} from '@influxdata/clockface'

// Types
import {ViewType} from 'src/types/cells'
import {AppState} from 'src/types/stores'

// Constants
import {VIS_GRAPHICS} from 'src/shared/constants/vis'

// Actions
import {setViewType} from 'src/timeMachine/actions'

const ViewTypeDropdown: FunctionComponent = () => {
  const dispatch = useDispatch()
  const viewType = useSelector((state: AppState) => {
    const properties = state.timeMachine.viewProperties
    return properties.type
  })

  const onSetViewType = (viewType: ViewType) => {
    dispatch(setViewType(viewType))
  }

  const getViewTypeGraphic = (viewType: ViewType) => {
    const {graphic, name} = VIS_GRAPHICS.find(
      graphic => graphic.type === viewType
    )

    return (
      <>
        <div className={'view-type-dropdown--graphic'}>{graphic}</div>
        <div className={'view-type-dropdown--name'}>{name}</div>
      </>
    )
  }

  const dropdownItems = () => {
    return VIS_GRAPHICS.filter(g => {
      if (g.type === 'mosaic') {
        return false
      }

      return g.type !== 'band'
    }).map(g => (
      <Dropdown.Item
        key={`view-type--${g.type}`}
        id={`${g.type}`}
        value={g.type}
        onClick={onSetViewType}
        selected={`${g.type}` === viewType}
      >
        {getViewTypeGraphic(g.type)}
      </Dropdown.Item>
    ))
  }

  return (
    <Dropdown
      style={{width: '215px'}}
      className={'view-type-dropdown'}
      button={(active, onClick) => (
        <Dropdown.Button
          active={active}
          onClick={onClick}
          status={ComponentStatus.Valid}
        >
          {getViewTypeGraphic(viewType)}
        </Dropdown.Button>
      )}
      menu={onCollapse => (
        <Dropdown.Menu onCollapse={onCollapse} theme={DropdownMenuTheme.Onyx}>
          {dropdownItems()}
        </Dropdown.Menu>
      )}
    />
  )
}

export default ViewTypeDropdown
