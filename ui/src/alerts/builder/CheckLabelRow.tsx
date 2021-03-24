// Libraries
import React from 'react'

// Types
import {CheckLabel} from '../../types/Check'
import {
  ComponentColor,
  ComponentSize,
  DismissButton,
  FlexBox,
  FlexDirection,
  Input,
  Panel,
  TextBlock,
} from '@influxdata/clockface'

interface Props {
  index: number
  label: CheckLabel
  handleChangeLabelRow: (index: number, k: string, v: string) => void
  handleRemoveTagRow: (i: number) => void
}

const CheckLabelRow: React.FC<Props> = props => {
  const {index, label, handleChangeLabelRow, handleRemoveTagRow} = props

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // handleChangeLabelRow()
    console.log('onchange', index, e.target.name, e.target.value)
  }

  return (
    <Panel testID={'check-label'} className={'alert-builder--tag-row'}>
      <DismissButton
        color={ComponentColor.Default}
        onClick={() => {
          handleRemoveTagRow(index)
        }}
      />

      <Panel.Body size={ComponentSize.ExtraSmall}>
        <FlexBox direction={FlexDirection.Row} margin={ComponentSize.Small}>
          <FlexBox.Child grow={1}>
            <Input
              testID={'label-value-key--input'}
              placeholder={'Key'}
              value={label.key}
              name={'key'}
              onChange={handleChange}
            />
          </FlexBox.Child>

          <FlexBox.Child grow={0} basis={20}>
            <TextBlock text={'='} />
          </FlexBox.Child>

          <FlexBox.Child grow={1}>
            <Input
              testID={'label-rule-key--input'}
              placeholder={'Value'}
              name={'value'}
              onChange={handleChange}
            />
          </FlexBox.Child>
        </FlexBox>
      </Panel.Body>
    </Panel>
  )
}

export default CheckLabelRow
