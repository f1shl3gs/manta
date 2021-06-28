// Libraries
import React from 'react'

// Component
import {ButtonShape, SelectGroup} from '@influxdata/clockface'
import {useViewType, ViewType} from './useProfile'

const options = [ViewType.Table, ViewType.Both, ViewType.FlameGraph]

const FlameGraphControlBar: React.FC = () => {
  const {viewType, setViewType} = useViewType()

  return (
    <SelectGroup shape={ButtonShape.StretchToFit}>
      {options.map(key => (
        <SelectGroup.Option
          key={key}
          value={key}
          active={key === viewType}
          id={key}
          onClick={() => setViewType(key)}
        >
          {key}
        </SelectGroup.Option>
      ))}
    </SelectGroup>
  )
}

export default FlameGraphControlBar
