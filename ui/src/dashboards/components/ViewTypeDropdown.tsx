import React, {FunctionComponent, useState} from 'react'
import {
  ComponentStatus,
  Dropdown,
  DropdownMenuTheme,
} from '@influxdata/clockface'
import {
  SingleStatViewProperties,
  ViewProperties,
  ViewType,
} from 'src/types/Dashboard'
import {VIS_GRAPHICS} from 'src/constants/vis'

const defaultGaugeViewProperties = {
  type: 'gauge',
  queries: [],
  prefix: '',
  tickPrefix: '',
  suffix: '',
  tickSuffix: '',
  colors: [
    {
      id: '0',
      type: 'min',
      hex: '#00C9FF',
      name: 'laser',
      value: 0,
    },
    {
      id: '608289e1-a646-47c8-8602-4b2dd14c5751',
      type: 'threshold',
      hex: '#FFB94A',
      name: 'pineapple',
      value: 20,
    },
    {
      id: 'a2ae0345-1189-43de-8cf1-8aaa441bebd9',
      type: 'threshold',
      hex: '#DC4E58',
      name: 'fire',
      value: 85,
    },
    {
      id: '1',
      type: 'max',
      hex: '#9394FF',
      name: 'comet',
      value: 100,
    },
  ],
  decimalPlaces: {
    isEnforced: true,
    digits: 2,
  },
  note: '',
  showNoteWhenEmpty: false,
}

const defaultSingleStatViewProperties: SingleStatViewProperties = {
  type: 'single-stat',
  queries: [],
  colors: [],
  prefix: '',
  suffix: '',
  tickPrefix: '',
  tickSuffix: '',
  legend: {},
  showNoteWhenEmpty: true,
  note: '',
  decimalPlaces: {
    isEnforced: true,
    digits: 2,
  },
}

interface Props {
  setViewProperties: (ViewProperties) => void
}

const ViewTypeDropdown: FunctionComponent<Props> = ({setViewProperties}) => {
  const [viewType, setViewType] = useState<ViewType>('xy')

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
        onClick={value => {
          setViewType(value)
          if (value === 'gauge') {
            setViewProperties(defaultGaugeViewProperties as ViewProperties)
          } else if (value === 'single-stat') {
            setViewProperties(defaultSingleStatViewProperties)
          }
        }}
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
