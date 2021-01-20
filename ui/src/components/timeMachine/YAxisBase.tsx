// Libraries
import React from 'react'

// Components
import {
  ButtonShape,
  Columns,
  Form,
  Grid,
  SelectGroup,
} from '@influxdata/clockface'

// Types
import {AXES_SCALE_OPTIONS} from 'constants/cell'

interface Props {
  base: string
  onSetYAxisBase: (base: string) => void
}

const {BASE_2, BASE_10} = AXES_SCALE_OPTIONS

const options: {
  id: string
  name: string
  value: string
  titleText: string
  displayText: string
}[] = [
  {
    id: 'y-values-format-tab--raw',
    name: 'y-values-format',
    titleText: 'Do not format values using a unit prefix',
    value: '',
    displayText: 'None',
  },
  {
    id: 'y-values-format-tab--si',
    name: 'y-values-format',
    titleText: 'Format values using an International System of Units prefix',
    value: BASE_10 as string,
    displayText: 'SI',
  },
  {
    id: 'y-values-format-tab--binary',
    name: 'y-values-format',
    titleText:
      'Format values using a binary unit prefix (for formatting bits or bytes)',
    value: BASE_2 as string,
    displayText: 'Binary',
  },
]

const YAxisBase: React.FC<Props> = (props) => {
  const {base, onSetYAxisBase} = props

  return (
    <Grid.Column widthXS={Columns.Twelve}>
      <Form.Element label={'Y-Value Unit Prefix'}>
        <SelectGroup shape={ButtonShape.StretchToFit}>
          {options.map((item) => (
            <SelectGroup.Option
              key={item.id}
              name={item.name}
              id={item.id}
              value={item.value}
              active={base === item.value}
              titleText={item.titleText}
              onClick={onSetYAxisBase}
            >
              {item.displayText}
            </SelectGroup.Option>
          ))}
        </SelectGroup>
      </Form.Element>
    </Grid.Column>
  )
}

export default YAxisBase
