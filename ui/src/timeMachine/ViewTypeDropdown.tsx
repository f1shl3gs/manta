import React, {FunctionComponent} from 'react'
import {
  ComponentStatus,
  Dropdown,
  DropdownMenuTheme,
} from '@influxdata/clockface'
import {ViewType} from 'src/types/cells'
import {VIS_GRAPHICS} from 'src/constants/vis'
import {useDispatch, useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'

const ViewTypeDropdown: FunctionComponent = () => {
  const dispatch = useDispatch()
  const viewType = useSelector((state: AppState) => {
    const properties = state.timeMachine.viewProperties
    return properties.type
  })

  const setViewType = (viewType: ViewType) => {
    dispatch(setViewType(viewType))
  }

  const getViewTypeGraphic = (viewType: ViewType) => {
    // @ts-ignore
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
        onClick={setViewType}
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
